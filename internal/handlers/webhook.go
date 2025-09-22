package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"novinhub-webhook/internal/config"
	"novinhub-webhook/internal/models"
	"novinhub-webhook/internal/services"
	"novinhub-webhook/internal/utils"
	"novinhub-webhook/pkg/logger"
)

// SMSCache represents a simple cache entry for SMS deduplication
type SMSCache struct {
	SentAt time.Time
	Phone  string
	UserID string
}

// WebhookHandler handles incoming webhook requests
type WebhookHandler struct {
	logger     *logger.Logger
	smsService *services.SMSService
	smsCache   map[string]SMSCache // key: phone_userID, value: cache entry
	cacheMutex sync.RWMutex        // mutex for thread-safe cache operations
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(logger *logger.Logger, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{
		logger:     logger,
		smsService: services.NewSMSService(logger, cfg),
		smsCache:   make(map[string]SMSCache),
		cacheMutex: sync.RWMutex{},
	}
}

// shouldSendSMS checks if SMS should be sent (daily limit logic)
func (h *WebhookHandler) shouldSendSMS(phone, userID string) bool {
	cacheKey := fmt.Sprintf("%s_%s", phone, userID)

	h.cacheMutex.RLock()
	cached, exists := h.smsCache[cacheKey]
	h.cacheMutex.RUnlock()

	if !exists {
		return true // No previous SMS sent
	}

	// Check if it's a new day (daily reset)
	now := time.Now()
	lastSMSDate := cached.SentAt.Truncate(24 * time.Hour) // Get date only (remove time)
	todayDate := now.Truncate(24 * time.Hour)

	if lastSMSDate.Before(todayDate) {
		// It's a new day, allow SMS
		h.logger.Info("✅ NEW DAY - SMS ALLOWED",
			"phone", phone,
			"user_id", userID,
			"last_sms_date", lastSMSDate.Format("2006-01-02"),
			"today_date", todayDate.Format("2006-01-02"))
		return true
	}

	// Same day - block SMS
	timeSinceLast := time.Since(cached.SentAt)
	timeUntilNextDay := time.Until(todayDate.Add(24 * time.Hour))

	h.logger.Warn("🚫 SMS BLOCKED - ALREADY SENT TODAY",
		"phone", phone,
		"user_id", userID,
		"last_sms_time", cached.SentAt.Format("2006-01-02 15:04:05"),
		"time_since_last", timeSinceLast.String(),
		"next_allowed_in", timeUntilNextDay.String())

	return false
}

// markSMSSent records that SMS was sent for deduplication
func (h *WebhookHandler) markSMSSent(phone, userID string) {
	cacheKey := fmt.Sprintf("%s_%s", phone, userID)

	h.cacheMutex.Lock()
	h.smsCache[cacheKey] = SMSCache{
		SentAt: time.Now(),
		Phone:  phone,
		UserID: userID,
	}
	h.cacheMutex.Unlock()

	h.logger.Info("📝 SMS CACHE UPDATED",
		"phone", phone,
		"user_id", userID,
		"cache_key", cacheKey)

	// Clean old cache entries (older than 3 days)
	go h.cleanupOldCache()
}

// cleanupOldCache removes cache entries older than 3 days
func (h *WebhookHandler) cleanupOldCache() {
	h.cacheMutex.Lock()
	defer h.cacheMutex.Unlock()

	cutoffTime := time.Now().Add(-3 * 24 * time.Hour) // 3 days ago
	removedCount := 0

	for key, cached := range h.smsCache {
		if cached.SentAt.Before(cutoffTime) {
			delete(h.smsCache, key)
			removedCount++
		}
	}

	if removedCount > 0 {
		h.logger.Info("🧹 CACHE CLEANUP COMPLETED",
			"removed_entries", removedCount,
			"remaining_entries", len(h.smsCache),
			"cutoff_date", cutoffTime.Format("2006-01-02"))
	}
}

// HandleWebhook processes incoming webhook events
func (h *WebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	// Only accept POST requests
	if r.Method != http.MethodPost {
		h.logger.Warn("Invalid method received", "method", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Read the raw body first for logging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read request body", "error", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Log the complete raw request
	h.logger.Info("Raw webhook request received",
		"method", r.Method,
		"url", r.URL.String(),
		"headers", r.Header,
		"raw_body", string(body))

	// Parse the webhook event
	var event models.WebhookEvent
	if err := json.Unmarshal(body, &event); err != nil {
		h.logger.Error("Failed to decode webhook payload", "error", err)
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Log the received event with raw payload
	rawPayload, _ := json.Marshal(event.Payload)
	h.logger.Info("Webhook event received",
		"type", event.Type,
		"user_id", event.UserID.String(),
		"timestamp", time.Now().UTC(),
		"raw_payload", string(rawPayload))

	// Process different event types
	switch event.Type {
	case "message_created":
		h.handleMessageCreated(event)
	case "comment_created":
		h.handleCommentCreated(event)
	case "autoform_completed":
		h.handleAutoformCompleted(event)
	case "leed_created":
		h.handleLeadCreated(event)
	case "revalidate":
		h.handleRevalidate(event)
	default:
		h.logger.Warn("Unknown event type received", "type", event.Type)
	}

	// Return 200 OK as required by NovinHub
	w.WriteHeader(http.StatusOK)
	response := models.WebhookResponse{
		Status:  "success",
		Message: "Webhook processed successfully",
	}
	json.NewEncoder(w).Encode(response)
}

// handleMessageCreated processes message_created events
func (h *WebhookHandler) handleMessageCreated(event models.WebhookEvent) {
	h.logger.Info("Processing message_created event", "user_id", event.UserID.String())

	// Parse the message payload
	messageBytes, err := json.Marshal(event.Payload)
	if err != nil {
		h.logger.Error("Failed to marshal message payload", "error", err)
		return
	}

	var message models.Message
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		h.logger.Error("Failed to unmarshal message", "error", err)
		return
	}

	h.logger.Info("Message details",
		"message_id", message.ID,
		"text", message.Text,
		"content", message.Content,
		"account", message.Account,
		"social_user", message.SocialUser)

	// Extract phone numbers from message text for logging only
	// Note: We don't send SMS here to avoid duplicates with leed_created event
	if message.Text != "" {
		phoneNumbers := utils.ExtractIranianPhoneNumbers(message.Text)
		if len(phoneNumbers) > 0 {
			h.logger.Info("📱 PHONE NUMBER DETECTED IN DIRECT MESSAGE! 📱",
				"user_id", event.UserID.String(),
				"message_id", message.ID,
				"phone_numbers", phoneNumbers,
				"message_text", message.Text,
				"note", "SMS will be sent via leed_created event to avoid duplicates")

			// Log detected numbers but don't send SMS (wait for leed_created)
			for _, phone := range phoneNumbers {
				if utils.IsValidIranianPhone(phone) {
					h.logger.Info("💡 PHONE DETECTED - WAITING FOR LEAD EVENT",
						"phone", phone,
						"user_id", event.UserID.String(),
						"message_id", message.ID,
						"status", "awaiting_lead_event")
				}
			}
		}
	}

	// Add your business logic here for handling new messages
	// For example: save to database, send notifications, etc.
}

// handleCommentCreated processes comment_created events
func (h *WebhookHandler) handleCommentCreated(event models.WebhookEvent) {
	h.logger.Info("Processing comment_created event", "user_id", event.UserID.String())

	// Parse the comment payload
	commentBytes, err := json.Marshal(event.Payload)
	if err != nil {
		h.logger.Error("Failed to marshal comment payload", "error", err)
		return
	}

	var comment models.Comment
	if err := json.Unmarshal(commentBytes, &comment); err != nil {
		h.logger.Error("Failed to unmarshal comment", "error", err)
		return
	}

	h.logger.Info("Comment details",
		"comment_id", comment.ID,
		"content", comment.Content,
		"account", comment.Account,
		"social_user", comment.SocialUser,
		"account_post", comment.AccountPost)

	// Add your business logic here for handling new comments
}

// handleAutoformCompleted processes autoform_completed events
func (h *WebhookHandler) handleAutoformCompleted(event models.WebhookEvent) {
	h.logger.Info("Processing autoform_completed event", "user_id", event.UserID.String())

	// Parse the form response payload
	formBytes, err := json.Marshal(event.Payload)
	if err != nil {
		h.logger.Error("Failed to marshal form response payload", "error", err)
		return
	}

	var formResponse models.AutomationFormResponse
	if err := json.Unmarshal(formBytes, &formResponse); err != nil {
		h.logger.Error("Failed to unmarshal form response", "error", err)
		return
	}

	h.logger.Info("Form response details",
		"form_id", formResponse.ID,
		"messages", formResponse.Messages,
		"social_user", formResponse.SocialUser)

	// Add your business logic here for handling completed forms
}

// handleLeadCreated processes leed_created events
func (h *WebhookHandler) handleLeadCreated(event models.WebhookEvent) {
	h.logger.Info("Processing leed_created event", "user_id", event.UserID.String())

	// Parse the lead payload
	leadBytes, err := json.Marshal(event.Payload)
	if err != nil {
		h.logger.Error("Failed to marshal lead payload", "error", err)
		return
	}

	var lead models.Lead
	if err := json.Unmarshal(leadBytes, &lead); err != nil {
		h.logger.Error("Failed to unmarshal lead", "error", err)
		return
	}

	h.logger.Info("Lead details",
		"lead_id", lead.ID,
		"type", lead.Type,
		"value", lead.Value,
		"message_id", lead.MessageID,
		"social_user", lead.SocialUser)

	// Process phone number leads specifically
	if lead.Type == "number" && lead.Value != "" {
		// Validate if it's a valid Iranian phone number
		if utils.IsValidIranianPhone(lead.Value) {
			h.logger.Warn("🎯 LEAD WITH VALID PHONE NUMBER DETECTED! 🎯",
				"phone", lead.Value,
				"lead_id", lead.ID,
				"user_id", event.UserID.String(),
				"message_id", lead.MessageID)

			// Check if we should send SMS (deduplication)
			if h.shouldSendSMS(lead.Value, event.UserID.String()) {
				// Call SMS service to send pattern-based SMS
				err := h.smsService.SendSMSWithPattern(
					lead.Value,
					event.UserID.String(),
				)

				if err != nil {
					h.logger.Error("Failed to send SMS for lead",
						"error", err,
						"phone", lead.Value,
						"lead_id", lead.ID)
				} else {
					// Mark SMS as sent to prevent duplicates
					h.markSMSSent(lead.Value, event.UserID.String())

					h.logger.Info("✅ SMS PROCESSING COMPLETED FOR LEAD",
						"phone", lead.Value,
						"lead_id", lead.ID,
						"status", "sent_and_cached")
				}
			} else {
				h.logger.Info("⏭️ SMS SKIPPED - RECENTLY SENT",
					"phone", lead.Value,
					"lead_id", lead.ID,
					"user_id", event.UserID.String())
			}
		} else {
			h.logger.Warn("Invalid phone number in lead",
				"phone", lead.Value,
				"lead_id", lead.ID)
		}
	}

	// Add your business logic here for handling new leads
}

// handleRevalidate processes revalidate events
func (h *WebhookHandler) handleRevalidate(event models.WebhookEvent) {
	h.logger.Info("Processing revalidate event", "user_id", event.UserID.String())

	// Add your business logic here for handling revalidation
	// This is typically used to verify webhook authenticity
}
