package main

import (
	"log"

	"novinhub-webhook/internal/config"
	"novinhub-webhook/internal/server"
	"novinhub-webhook/pkg/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	// Initialize logger
	logger := logger.New()

	// Create server
	srv := server.New(cfg, logger)

	// Start server
	if err := srv.Start(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
