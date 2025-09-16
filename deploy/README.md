# Deployment Guide for NovinHub Webhook

This directory contains all the necessary files and scripts for deploying the NovinHub Webhook service to your Ubuntu VPS with domain `asllmarket.org`.

## 📁 Directory Structure

```
deploy/
├── nginx/
│   └── novinhub-webhook.conf          # Nginx configuration
├── supervisor/
│   └── novinhub-webhook.conf          # Supervisor configuration
└── scripts/
    ├── setup-vps.sh                   # Initial VPS setup
    ├── deploy.sh                      # Application deployment
    └── update.sh                      # Application update
```

## 🚀 Quick Deployment

### Step 1: Initial VPS Setup

SSH into your Ubuntu VPS and run:

```bash
# Download and run the setup script
wget https://raw.githubusercontent.com/your-repo/novinhub-webhook/main/deploy/scripts/setup-vps.sh
chmod +x setup-vps.sh
sudo ./setup-vps.sh
```

This script will:
- Install Go 1.21.5, Nginx, Supervisor, Certbot
- Configure UFW firewall and Fail2ban
- Set up SSL certificate for `asllmarket.org`
- Create application user `novinhub`
- Configure Nginx reverse proxy
- Set up log rotation

### Step 2: Deploy Application

On your local machine:

```bash
# Upload files to VPS
scp -r . user@asllmarket.org:/tmp/novinhub-webhook/

# SSH into VPS
ssh user@asllmarket.org

# Navigate to uploaded directory
cd /tmp/novinhub-webhook

# Run deployment
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

### Step 3: Verify Deployment

```bash
# Check service status
sudo supervisorctl status novinhub-webhook

# Check logs
sudo supervisorctl tail novinhub-webhook

# Test health endpoint
curl https://asllmarket.org/health

# Test webhook endpoint
curl -X POST https://asllmarket.org/webhook \
  -H "Content-Type: application/json" \
  -d '{"type":"test","user_id":"123","payload":{}}'
```

## 🔧 Configuration

### Production Configuration

The production configuration is located at `/opt/novinhub-webhook/config/config.yaml`:

```yaml
server:
  port: 8080
  host: "127.0.0.1"  # Localhost for nginx proxy

logging:
  level: "info"
  format: "json"
  output: "file"
  file_path: "/var/log/novinhub-webhook/webhook.log"

environment:
  mode: "production"
  debug: false
```

### Nginx Configuration

Nginx is configured to:
- Handle SSL termination
- Rate limiting (10 requests/second)
- Security headers
- Proxy to localhost:8080
- Log all requests

### Supervisor Configuration

Supervisor manages the application process:
- Auto-restart on failure
- Log rotation
- Process monitoring
- Environment variables

## 📊 Monitoring

### Service Management

```bash
# Start service
sudo supervisorctl start novinhub-webhook

# Stop service
sudo supervisorctl stop novinhub-webhook

# Restart service
sudo supervisorctl restart novinhub-webhook

# Check status
sudo supervisorctl status novinhub-webhook

# View logs
sudo supervisorctl tail novinhub-webhook
```

### Log Files

- Application logs: `/var/log/novinhub-webhook/webhook.log`
- Supervisor logs: `/var/log/novinhub-webhook/supervisor.log`
- Nginx access logs: `/var/log/nginx/novinhub-webhook.access.log`
- Nginx error logs: `/var/log/nginx/novinhub-webhook.error.log`

### Health Checks

- Health endpoint: `https://asllmarket.org/health`
- Webhook endpoint: `https://asllmarket.org/webhook`

## 🔄 Updates

To update the application:

```bash
# Upload new files
scp -r . user@asllmarket.org:/tmp/novinhub-webhook/

# SSH into VPS
ssh user@asllmarket.org

# Navigate to directory
cd /tmp/novinhub-webhook

# Run update script
chmod +x deploy/scripts/update.sh
./deploy/scripts/update.sh
```

## 🔒 Security

The deployment includes several security measures:

- **SSL/TLS**: Automatic SSL certificate with Let's Encrypt
- **Firewall**: UFW configured with minimal open ports
- **Fail2ban**: Protection against brute force attacks
- **Rate Limiting**: Nginx rate limiting (10 req/s)
- **Security Headers**: HSTS, XSS protection, etc.
- **Process Isolation**: Application runs as non-root user

## 🐛 Troubleshooting

### Common Issues

1. **Service won't start**:
   ```bash
   sudo supervisorctl tail novinhub-webhook
   ```

2. **Nginx errors**:
   ```bash
   sudo nginx -t
   sudo tail -f /var/log/nginx/error.log
   ```

3. **SSL certificate issues**:
   ```bash
   sudo certbot certificates
   sudo certbot renew --dry-run
   ```

4. **Permission issues**:
   ```bash
   sudo chown -R novinhub:novinhub /opt/novinhub-webhook
   sudo chown -R novinhub:novinhub /var/log/novinhub-webhook
   ```

### Log Analysis

```bash
# Application logs
tail -f /var/log/novinhub-webhook/webhook.log

# Nginx logs
tail -f /var/log/nginx/novinhub-webhook.access.log
tail -f /var/log/nginx/novinhub-webhook.error.log

# System logs
journalctl -u nginx -f
```

## 📈 Performance

### Monitoring Commands

```bash
# Check memory usage
ps aux | grep webhook

# Check disk usage
df -h

# Check network connections
netstat -tlnp | grep :8080

# Check nginx status
systemctl status nginx
```

### Optimization

- Nginx is configured with keepalive connections
- Application logs are rotated automatically
- Rate limiting prevents abuse
- SSL is optimized for performance

## 🔄 Backup

### Important Files to Backup

- `/opt/novinhub-webhook/config/config.yaml`
- `/etc/nginx/sites-available/novinhub-webhook.conf`
- `/etc/supervisor/conf.d/novinhub-webhook.conf`
- SSL certificates in `/etc/letsencrypt/`

### Backup Script

```bash
#!/bin/bash
# Create backup
tar -czf novinhub-webhook-backup-$(date +%Y%m%d).tar.gz \
  /opt/novinhub-webhook/config/ \
  /etc/nginx/sites-available/novinhub-webhook.conf \
  /etc/supervisor/conf.d/novinhub-webhook.conf \
  /etc/letsencrypt/live/asllmarket.org/
```
