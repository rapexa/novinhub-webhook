package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"novinhub-webhook/internal/models"
	"novinhub-webhook/pkg/logger"
)

// WebhookHandler handles incoming webhook requests
type WebhookHandler struct {
	logger *logger.Logger
}

// NewWebhookHandler creates a new webhook handler
func NewWebhookHandler(logger *logger.Logger) *WebhookHandler {
	return &WebhookHandler{
		logger: logger,
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
		"content", message.Content,
		"account", message.Account,
		"social_user", message.SocialUser)

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
		"phone", lead.Phone,
		"messages", lead.Messages,
		"social_user", lead.SocialUser)

	// Add your business logic here for handling new leads
}

// handleRevalidate processes revalidate events
func (h *WebhookHandler) handleRevalidate(event models.WebhookEvent) {
	h.logger.Info("Processing revalidate event", "user_id", event.UserID.String())

	// Add your business logic here for handling revalidation
	// This is typically used to verify webhook authenticity
}
