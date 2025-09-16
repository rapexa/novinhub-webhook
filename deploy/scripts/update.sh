#!/bin/bash

# Update Script for NovinHub Webhook
# Run this script to update the application

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="novinhub-webhook"
APP_DIR="/opt/$APP_NAME"
SERVICE_NAME="novinhub-webhook"

# Functions
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   log_error "This script should not be run as root"
   exit 1
fi

log_info "Updating $APP_NAME..."

# Stop the service
log_info "Stopping service..."
sudo supervisorctl stop $SERVICE_NAME

# Build the application
log_info "Building application..."
go build -o build/$APP_NAME cmd/server/main.go

# Create backup
log_info "Creating backup..."
sudo cp $APP_DIR/bin/$APP_NAME $APP_DIR/bin/$APP_NAME.backup.$(date +%Y%m%d_%H%M%S)

# Copy new binary
log_info "Installing new binary..."
sudo cp build/$APP_NAME $APP_DIR/bin/
sudo chmod +x $APP_DIR/bin/$APP_NAME
sudo chown novinhub:novinhub $APP_DIR/bin/$APP_NAME

# Update configuration if needed
if [ -f "config.production.yaml" ]; then
    log_info "Updating configuration..."
    sudo cp config.production.yaml $APP_DIR/config/
    sudo chown novinhub:novinhub $APP_DIR/config/config.production.yaml
fi

# Start the service
log_info "Starting service..."
sudo supervisorctl start $SERVICE_NAME

# Check service status
log_info "Checking service status..."
sudo supervisorctl status $SERVICE_NAME

# Test the service
log_info "Testing service..."
sleep 5
if curl -f http://localhost:8080/health > /dev/null 2>&1; then
    log_info "Service is running correctly!"
else
    log_error "Service failed to start properly. Check logs:"
    log_error "sudo supervisorctl tail $SERVICE_NAME"
    exit 1
fi

log_info "Update completed successfully!"
