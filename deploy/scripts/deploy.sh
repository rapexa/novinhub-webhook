#!/bin/bash

# NovinHub Webhook Deployment Script for Ubuntu VPS
# Domain: asllmarket.org

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
APP_NAME="novinhub-webhook"
APP_USER="novinhub"
APP_DIR="/opt/$APP_NAME"
SERVICE_NAME="novinhub-webhook"
DOMAIN="asllmarket.org"

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

# Check if user exists
if ! id "$APP_USER" &>/dev/null; then
    log_info "Creating user $APP_USER..."
    sudo useradd -r -s /bin/false -d $APP_DIR $APP_USER
fi

# Create application directory
log_info "Creating application directory..."
sudo mkdir -p $APP_DIR/{bin,config,logs}
sudo chown -R $APP_USER:$APP_USER $APP_DIR

# Create log directory
sudo mkdir -p /var/log/$APP_NAME
sudo chown -R $APP_USER:$APP_USER /var/log/$APP_NAME

# Build the application
log_info "Building application..."
go build -o novinhub-webhook cmd/server/main.go

# Copy binary
log_info "Installing binary..."
sudo cp build/$APP_NAME $APP_DIR/bin/
sudo chmod +x $APP_DIR/bin/$APP_NAME
sudo chown $APP_USER:$APP_USER $APP_DIR/bin/$APP_NAME

# Copy configuration
log_info "Installing configuration..."
sudo cp config.production.yaml $APP_DIR/config/
sudo chown $APP_USER:$APP_USER $APP_DIR/config/config.production.yaml

# Install supervisor configuration
log_info "Installing supervisor configuration..."
sudo cp deploy/supervisor/$SERVICE_NAME.conf /etc/supervisor/conf.d/
sudo supervisorctl reread
sudo supervisorctl update

# Install nginx configuration
log_info "Installing nginx configuration..."
sudo cp deploy/nginx/$SERVICE_NAME.conf /etc/nginx/sites-available/
sudo ln -sf /etc/nginx/sites-available/$SERVICE_NAME.conf /etc/nginx/sites-enabled/

# Test nginx configuration
log_info "Testing nginx configuration..."
sudo nginx -t

# Reload nginx
log_info "Reloading nginx..."
sudo systemctl reload nginx

# Start the service
log_info "Starting $SERVICE_NAME service..."
sudo supervisorctl start $SERVICE_NAME

# Check service status
log_info "Checking service status..."
sudo supervisorctl status $SERVICE_NAME

log_info "Deployment completed successfully!"
log_info "Webhook URL: http://$DOMAIN/webhook"
log_info "Health check: http://$DOMAIN/health"
log_info "Logs: /var/log/$APP_NAME/"
log_info "Note: SSL will be added later for HTTPS support"
