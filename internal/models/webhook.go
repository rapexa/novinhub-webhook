package models

// WebhookEvent represents the structure of webhook data from NovinHub
type WebhookEvent struct {
	Type    string      `json:"type"`
	UserID  string      `json:"user_id"`
	Payload interface{} `json:"payload"`
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
