package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"novinhub-webhook/internal/models"
	"novinhub-webhook/pkg/logger"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	logger *logger.Logger
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(logger *logger.Logger) *HealthHandler {
	return &HealthHandler{
		logger: logger,
	}
}

// HealthCheck endpoint for monitoring
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := models.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}
