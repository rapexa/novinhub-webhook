package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// FlexibleUserID can handle both string and numeric user IDs
type FlexibleUserID string

// UnmarshalJSON implements custom unmarshaling for FlexibleUserID
func (f *FlexibleUserID) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		*f = FlexibleUserID(str)
		return nil
	}

	// If that fails, try as number
	var num float64
	if err := json.Unmarshal(data, &num); err == nil {
		*f = FlexibleUserID(strconv.FormatFloat(num, 'f', -1, 64))
		return nil
	}

	return fmt.Errorf("user_id must be either string or number")
}

// String returns the string representation
func (f FlexibleUserID) String() string {
	return string(f)
}

// WebhookEvent represents the structure of webhook data from NovinHub
type WebhookEvent struct {
	Type    string         `json:"type"`
	UserID  FlexibleUserID `json:"user_id"`
	Payload interface{}    `json:"payload"`
}

// Message represents a message created event
type Message struct {
	ID         string      `json:"id"`
	Content    string      `json:"content"`
	Account    interface{} `json:"account"`
	SocialUser interface{} `json:"socialUser"`
}

// Comment represents a comment created event
type Comment struct {
	ID          string      `json:"id"`
	Content     string      `json:"content"`
	Account     interface{} `json:"account"`
	SocialUser  interface{} `json:"socialUser"`
	AccountPost interface{} `json:"accountPost"`
}

// AutomationFormResponse represents an autoform completed event
type AutomationFormResponse struct {
	ID         string      `json:"id"`
	Messages   interface{} `json:"messages"`
	SocialUser interface{} `json:"socialUser"`
}

// Lead represents a lead created event
type Lead struct {
	ID         string      `json:"id"`
	Phone      string      `json:"phone"`
	Messages   interface{} `json:"messages"`
	SocialUser interface{} `json:"socialUser"`
}

// WebhookResponse represents the response sent back to NovinHub
type WebhookResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}
