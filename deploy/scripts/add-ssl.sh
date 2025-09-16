#!/bin/bash

# Add SSL Script for NovinHub Webhook
# Run this script to add SSL support to the existing HTTP configuration

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DOMAIN="asllmarket.org"
APP_NAME="novinhub-webhook"

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
if [[ $EUID -ne 0 ]]; then
   log_error "This script must be run as root"
   exit 1
fi

log_info "Adding SSL support to NovinHub Webhook..."

# Install certbot if not already installed
if ! command -v certbot &> /dev/null; then
    log_info "Installing certbot..."
    apt update
    apt install -y certbot python3-certbot-nginx
fi

# Get SSL certificate
log_info "Obtaining SSL certificate for $DOMAIN..."
certbot --nginx -d $DOMAIN -d www.$DOMAIN --non-interactive --agree-tos --email admin@$DOMAIN

# Setup SSL renewal
log_info "Setting up SSL renewal..."
(crontab -l 2>/dev/null; echo "0 12 * * * /usr/bin/certbot renew --quiet") | crontab -

# Test nginx configuration
log_info "Testing nginx configuration..."
nginx -t

# Reload nginx
log_info "Reloading nginx..."
systemctl reload nginx

log_info "SSL setup completed successfully!"
log_info "Webhook URL: https://$DOMAIN/webhook"
log_info "Health check: https://$DOMAIN/health"
log_info "SSL certificate will auto-renew via cron job"
