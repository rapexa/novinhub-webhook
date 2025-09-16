package server

import (
	"log"
	"net/http"
	"time"

	"novinhub-webhook/internal/config"
	"novinhub-webhook/internal/handlers"
	"novinhub-webhook/pkg/logger"

	"github.com/gorilla/mux"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	logger  *logger.Logger
	handler *handlers.WebhookHandler
	health  *handlers.HealthHandler
}

// New creates a new server instance
func New(cfg *config.Config, log *logger.Logger) *Server {
	webhookHandler := handlers.NewWebhookHandler(log)
	healthHandler := handlers.NewHealthHandler(log)

	return &Server{
		config:  cfg,
		logger:  log,
		handler: webhookHandler,
		health:  healthHandler,
	}
}

// SetupRoutes configures the HTTP routes
func (s *Server) SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Webhook endpoint
	router.HandleFunc("/webhook", s.handler.HandleWebhook).Methods("POST")

	// Health check endpoint
	router.HandleFunc("/health", s.health.HealthCheck).Methods("GET")

	// Add CORS middleware
	router.Use(s.corsMiddleware)

	return router
}

// corsMiddleware adds CORS headers
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	router := s.SetupRoutes()

	server := &http.Server{
		Addr:         s.config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	s.logger.Info("Starting webhook server", "address", s.config.GetServerAddress())
	s.logger.Info("Webhook endpoint available at: http://" + s.config.GetServerAddress() + "/webhook")
	s.logger.Info("Health check available at: http://" + s.config.GetServerAddress() + "/health")

	return server.ListenAndServe()
}

// StartWithGracefulShutdown starts the server with graceful shutdown
func (s *Server) StartWithGracefulShutdown() error {
	router := s.SetupRoutes()

	server := &http.Server{
		Addr:         s.config.GetServerAddress(),
		Handler:      router,
		ReadTimeout:  time.Duration(s.config.Server.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(s.config.Server.WriteTimeout) * time.Second,
	}

	s.logger.Info("Starting webhook server", "address", s.config.GetServerAddress())
	s.logger.Info("Webhook endpoint available at: http://" + s.config.GetServerAddress() + "/webhook")
	s.logger.Info("Health check available at: http://" + s.config.GetServerAddress() + "/health")

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	// This would require signal handling in a real application
	// For now, we'll just return the server
	return nil
}
