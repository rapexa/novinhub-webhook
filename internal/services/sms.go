package services

import (
	"fmt"

	"novinhub-webhook/internal/config"
	"novinhub-webhook/internal/utils"
	"novinhub-webhook/pkg/logger"
)

// SMSService handles SMS sending functionality
type SMSService struct {
	logger        *logger.Logger
	config        *config.Config
	ippanelClient *IPPanelClient
}

// NewSMSService creates a new SMS service instance
func NewSMSService(logger *logger.Logger, cfg *config.Config) *SMSService {
	var ippanelClient *IPPanelClient

	// Initialize IPPanel client if API key is provided
	if cfg.SMS.IPPanel.APIKey != "" {
		ippanelClient = NewIPPanelClient(cfg.SMS.IPPanel.APIKey)
		logger.Info("📡 IPPanel SMS Client initialized",
			"provider", cfg.SMS.Provider,
			"enabled", cfg.SMS.Enabled,
			"originator", cfg.SMS.IPPanel.Originator)
	} else {
		logger.Warn("⚠️ IPPanel API key not configured - SMS will be disabled")
	}

	return &SMSService{
		logger:        logger,
		config:        cfg,
		ippanelClient: ippanelClient,
	}
}

// SendSMSWithPattern sends an SMS with the current daily pattern to a phone number
func (s *SMSService) SendSMSWithPattern(phoneNumber string, userID string) error {
	// Get current pattern from config
	currentPattern, patternIndex, groupName := s.config.GetCurrentPatternInfo()

	// Prepare variables for logging
	var logCode string
	if userID == "" {
		logCode = "کاربر گرامی"
	} else {
		logCode = userID
	}

	s.logger.Info("📲 SMS SENDING INITIATED 📲",
		"phone", phoneNumber,
		"pattern", currentPattern,
		"pattern_group", groupName,
		"pattern_index", patternIndex,
		"user_id", userID,
		"enabled", s.config.SMS.Enabled,
		"pattern_variables", map[string]string{
			"code": logCode,
		})

	// Validate phone number
	if !utils.IsValidIranianPhone(phoneNumber) {
		return fmt.Errorf("invalid Iranian phone number: %s", phoneNumber)
	}

	// Check if SMS is enabled
	if !s.config.SMS.Enabled {
		s.logger.Warn("📵 SMS DISABLED - SKIPPING SEND",
			"phone", phoneNumber,
			"pattern", currentPattern,
			"status", "disabled_in_config")
		return nil
	}

	// Check if IPPanel client is configured
	if s.ippanelClient == nil {
		s.logger.Error("❌ SMS CLIENT NOT CONFIGURED",
			"phone", phoneNumber,
			"error", "IPPanel client is nil - check API key configuration")
		return fmt.Errorf("SMS client not configured")
	}

	// Check required configuration
	if s.config.SMS.IPPanel.Originator == "" || s.config.SMS.IPPanel.PatternCode == "" {
		s.logger.Error("❌ SMS CONFIGURATION INCOMPLETE",
			"phone", phoneNumber,
			"originator_configured", s.config.SMS.IPPanel.Originator != "",
			"pattern_configured", s.config.SMS.IPPanel.PatternCode != "")
		return fmt.Errorf("SMS configuration incomplete")
	}

	// Prepare pattern variables (customize as needed)
	// Only one variable: 'code' - if userID is empty, use "کاربر گرامی", otherwise use userID
	var code string
	if userID == "" {
		code = "کاربر گرامی"
	} else {
		code = userID
	}

	variables := map[string]string{
		"code": code,
	}

	// Send SMS using IPPanel with current pattern
	messageID, err := s.ippanelClient.SendPattern(
		currentPattern, // Use current pattern from pattern manager
		s.config.SMS.IPPanel.Originator,
		phoneNumber,
		variables,
	)

	if err != nil {
		s.logger.Error("❌ SMS SEND FAILED",
			"error", err,
			"phone", phoneNumber,
			"pattern", currentPattern)
		return fmt.Errorf("failed to send SMS: %v", err)
	}

	s.logger.Info("✅ SMS SENT SUCCESSFULLY",
		"phone", phoneNumber,
		"user_id", userID,
		"message_id", messageID,
		"pattern", currentPattern,
		"pattern_group", groupName,
		"pattern_index", patternIndex,
		"originator", s.config.SMS.IPPanel.Originator,
		"pattern_variables", map[string]string{
			"code": code,
		})

	return nil
}

// SendBulkSMS sends SMS to multiple phone numbers
func (s *SMSService) SendBulkSMS(phoneNumbers []string, message string) error {
	s.logger.Info("📱 BULK SMS INITIATED 📱",
		"count", len(phoneNumbers))

	// TODO: Implement bulk SMS sending logic
	s.logger.Warn("🚧 BULK SMS FUNCTION NOT IMPLEMENTED YET 🚧",
		"phone_count", len(phoneNumbers),
		"status", "TODO")

	return nil
}
