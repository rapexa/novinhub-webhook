#!/bin/bash

# VPS Setup Script for NovinHub Webhook
# Ubuntu 20.04/22.04 LTS

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
DOMAIN="asllmarket.org"
APP_USER="novinhub"
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

log_info "Setting up VPS for NovinHub Webhook..."

# Update system
log_info "Updating system packages..."
apt update && apt upgrade -y

# Install required packages
log_info "Installing required packages..."
apt install -y curl wget git nginx supervisor certbot python3-certbot-nginx ufw fail2ban

# Install Go
log_info "Installing Go..."
GO_VERSION="1.21.5"
wget https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz
tar -C /usr/local -xzf go${GO_VERSION}.linux-amd64.tar.gz
rm go${GO_VERSION}.linux-amd64.tar.gz

# Add Go to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> /etc/profile
echo 'export PATH=$PATH:/usr/local/go/bin' >> /root/.bashrc

# Create application user
log_info "Creating application user..."
useradd -r -s /bin/false -d /opt/$APP_NAME $APP_USER

# Create directories
log_info "Creating directories..."
mkdir -p /opt/$APP_NAME/{bin,config,logs}
mkdir -p /var/log/$APP_NAME
chown -R $APP_USER:$APP_USER /opt/$APP_NAME
chown -R $APP_USER:$APP_USER /var/log/$APP_NAME

# Configure firewall
log_info "Configuring firewall..."
ufw --force enable
ufw allow ssh
ufw allow 'Nginx Full'
ufw allow 80
ufw allow 443

# Configure fail2ban
log_info "Configuring fail2ban..."
cat > /etc/fail2ban/jail.local << EOF
[DEFAULT]
bantime = 3600
findtime = 600
maxretry = 3

[sshd]
enabled = true
port = ssh
logpath = /var/log/auth.log
maxretry = 3

[nginx-http-auth]
enabled = true
filter = nginx-http-auth
port = http,https
logpath = /var/log/nginx/error.log
maxretry = 3
EOF

systemctl enable fail2ban
systemctl start fail2ban

# Configure nginx
log_info "Configuring nginx..."
# Remove default site
rm -f /etc/nginx/sites-enabled/default

# Create nginx configuration (HTTP only for now)
cat > /etc/nginx/sites-available/$APP_NAME << EOF
# NovinHub Webhook Nginx Configuration
# Domain: $DOMAIN
# HTTP Only Configuration (SSL will be added later)

upstream novinhub_webhook {
    server 127.0.0.1:8080;
    keepalive 32;
}

# Rate limiting
limit_req_zone \$binary_remote_addr zone=webhook:10m rate=10r/s;

server {
    listen 80;
    server_name $DOMAIN www.$DOMAIN;
    
    # Security headers
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
    add_header X-XSS-Protection "1; mode=block";
    
    # Logging
    access_log /var/log/nginx/$APP_NAME.access.log;
    error_log /var/log/nginx/$APP_NAME.error.log;
    
    # Client settings
    client_max_body_size 1M;
    client_body_timeout 10s;
    client_header_timeout 10s;
    
    # Webhook endpoint
    location /webhook {
        # Rate limiting
        limit_req zone=webhook burst=20 nodelay;
        
        # Proxy settings
        proxy_pass http://novinhub_webhook;
        proxy_http_version 1.1;
        proxy_set_header Upgrade \$http_upgrade;
        proxy_set_header Connection 'upgrade';
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        proxy_cache_bypass \$http_upgrade;
        
        # Timeouts
        proxy_connect_timeout 5s;
        proxy_send_timeout 10s;
        proxy_read_timeout 10s;
        
        # Only allow POST requests
        limit_except POST {
            deny all;
        }
    }
    
    # Health check endpoint
    location /health {
        proxy_pass http://novinhub_webhook;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Real-IP \$remote_addr;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto \$scheme;
        
        # Allow GET requests only
        limit_except GET {
            deny all;
        }
    }
    
    # Block all other requests
    location / {
        return 404;
    }
}
EOF

# Enable site
ln -sf /etc/nginx/sites-available/$APP_NAME /etc/nginx/sites-enabled/

# Test nginx configuration
nginx -t

# Start nginx
systemctl enable nginx
systemctl start nginx

# Note: SSL certificate will be added later
log_info "Nginx configured for HTTP only. SSL will be added later."

# Configure supervisor
log_info "Configuring supervisor..."
cat > /etc/supervisor/conf.d/$APP_NAME.conf << EOF
[program:$APP_NAME]
command=/opt/$APP_NAME/bin/$APP_NAME
directory=/opt/$APP_NAME
user=$APP_USER
group=$APP_USER
autostart=true
autorestart=true
redirect_stderr=true
stdout_logfile=/var/log/$APP_NAME/supervisor.log
stdout_logfile_maxbytes=10MB
stdout_logfile_backups=5
environment=ENVIRONMENT=production
environment=CONFIG_PATH=/opt/$APP_NAME/config/config.yaml
EOF

# Reload supervisor
supervisorctl reread
supervisorctl update

log_info "VPS setup completed successfully!"
log_info "Next steps:"
log_info "1. Upload your application files to /opt/$APP_NAME/"
log_info "2. Run the deployment script"
log_info "3. Start the service with: supervisorctl start $APP_NAME"
log_info "4. Check status with: supervisorctl status $APP_NAME"
log_info "5. View logs with: tail -f /var/log/$APP_NAME/supervisor.log"
log_info "6. Webhook URL: http://$DOMAIN/webhook"
log_info "7. Health check: http://$DOMAIN/health"
log_info "Note: SSL will be added later for HTTPS support"
