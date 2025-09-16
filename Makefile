# Makefile for NovinHub Webhook Service

# Variables
BINARY_NAME=webhook
BUILD_DIR=build
MAIN_PATH=cmd/server/main.go

# Default target
.PHONY: all
all: build

# Build the application
.PHONY: build
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "Build complete: $(BUILD_DIR)/$(BINARY_NAME)"

# Run the application
.PHONY: run
run:
	@echo "Running $(BINARY_NAME)..."
	@go run $(MAIN_PATH)

# Run with hot reload (requires air)
.PHONY: dev
dev:
	@echo "Running in development mode..."
	@air

# Test the application
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@go clean

# Install dependencies
.PHONY: deps
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod tidy

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
.PHONY: lint
lint:
	@echo "Linting code..."
	@golangci-lint run

# Docker build
.PHONY: docker-build
docker-build:
	@echo "Building Docker image..."
	@docker build -t novinhub-webhook .

# Docker run
.PHONY: docker-run
docker-run:
	@echo "Running Docker container..."
	@docker run -p 8080:8080 novinhub-webhook

# Docker compose up
.PHONY: docker-up
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up

# Docker compose down
.PHONY: docker-down
docker-down:
	@echo "Stopping services with Docker Compose..."
	@docker-compose down

# Help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build        - Build the application"
	@echo "  run          - Run the application"
	@echo "  dev          - Run with hot reload (requires air)"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  deps         - Install dependencies"
	@echo "  fmt          - Format code"
	@echo "  lint         - Lint code"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-run   - Run Docker container"
	@echo "  docker-up    - Start with Docker Compose"
	@echo "  docker-down  - Stop Docker Compose services"
	@echo "  help         - Show this help message"
