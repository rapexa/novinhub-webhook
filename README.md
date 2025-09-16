# NovinHub Webhook Service

A well-structured Go-based webhook service to receive and process events from NovinHub platform.

## 🏗️ Project Structure

```
novinhub-webhook/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── handlers/
│   │   ├── webhook.go           # Webhook event handlers
│   │   └── health.go            # Health check handler
│   ├── models/
│   │   └── webhook.go           # Data models and structs
│   └── server/
│       └── server.go            # HTTP server setup
├── pkg/
│   └── logger/
│       └── logger.go            # Logging utilities
├── build/                       # Build artifacts
├── tmp/                         # Temporary files (air hot reload)
├── go.mod                       # Go module dependencies
├── go.sum                       # Go module checksums
├── Makefile                     # Build automation
├── Dockerfile                   # Docker configuration
├── docker-compose.yml           # Docker Compose setup
├── .air.toml                    # Hot reload configuration
├── .gitignore                   # Git ignore rules
└── README.md                    # This file
```

## ✨ Features

- ✅ **Clean Architecture**: Well-organized code with separation of concerns
- ✅ **All NovinHub Events**: Handles all webhook event types:
  - `message_created` - New direct messages
  - `comment_created` - New comments
  - `autoform_completed` - Completed smart forms
  - `leed_created` - New leads with phone numbers
  - `revalidate` - Webhook revalidation
- ✅ **Proper HTTP Response**: Returns 200 OK as required by NovinHub
- ✅ **Structured Logging**: JSON logs with context
- ✅ **Health Check**: Monitoring endpoint
- ✅ **CORS Support**: Cross-origin request handling
- ✅ **Error Handling**: Comprehensive error handling and validation
- ✅ **Hot Reload**: Development with automatic restarts
- ✅ **Docker Ready**: Containerized deployment
- ✅ **Makefile**: Easy build and run commands

## 🚀 Quick Start

### 1. Install Dependencies

```bash
make deps
# or
go mod tidy
```

### 2. Run the Service

```bash
# Development mode with hot reload
make dev

# Or run directly
make run

# Or with go
go run cmd/server/main.go
```

### 3. Webhook URLs

- **Webhook Endpoint**: `http://localhost:8080/webhook`
- **Health Check**: `http://localhost:8080/health`

## ⚙️ Configuration

### Configuration Files

The application uses Viper for configuration management with YAML files:

- `config.yaml` - Development configuration
- `config.production.yaml` - Production configuration

### Configuration Structure

```yaml
# Server configuration
server:
  port: 8080
  read_timeout: 10
  write_timeout: 10
  host: "0.0.0.0"

# Logging configuration
logging:
  level: "info"  # debug, info, warn, error
  format: "json" # json, text
  output: "stdout" # stdout, file
  file_path: "/var/log/novinhub-webhook/webhook.log"

# Webhook configuration
webhook:
  max_request_size: 1048576  # 1MB
  processing_timeout: 30
  enable_request_logging: true

# Security configuration
security:
  enable_cors: true
  allowed_origins: []
  rate_limit: 100

# Environment settings
environment:
  mode: "development"  # development, staging, production
  debug: false
```

### Environment Variables

You can override any configuration using environment variables:

- `SERVER_PORT` - Server port
- `SERVER_READ_TIMEOUT` - Read timeout in seconds
- `SERVER_WRITE_TIMEOUT` - Write timeout in seconds
- `LOGGING_LEVEL` - Log level
- `ENVIRONMENT_MODE` - Environment mode

### Webhook URL for NovinHub

When registering your webhook with NovinHub, use:
```
http://asllmarket.org/webhook
```

For local development with tools like ngrok:
```
https://your-ngrok-url.ngrok.io/webhook
```

**Note**: The production setup starts with HTTP only. SSL can be added later using the provided script.

## 📋 Available Commands

```bash
# Build the application
make build

# Run the application
make run

# Run with hot reload (development)
make dev

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Clean build artifacts
make clean

# Docker commands
make docker-build
make docker-run
make docker-up
make docker-down

# Show help
make help
```

## 🔧 Development

### Project Architecture

The project follows clean architecture principles:

- **`cmd/`**: Application entry points
- **`internal/`**: Private application code
  - **`config/`**: Configuration management
  - **`handlers/`**: HTTP request handlers
  - **`models/`**: Data models and structs
  - **`server/`**: HTTP server setup
- **`pkg/`**: Reusable packages
  - **`logger/`**: Logging utilities

### Adding Custom Logic

To add your own business logic for handling webhook events, modify the handler methods in `internal/handlers/webhook.go`:

- `handleMessageCreated()` - Process new messages
- `handleCommentCreated()` - Process new comments
- `handleAutoformCompleted()` - Process completed forms
- `handleLeadCreated()` - Process new leads
- `handleRevalidate()` - Handle revalidation

### Example: Adding Database Integration

```go
// In internal/handlers/webhook.go
func (h *WebhookHandler) handleMessageCreated(event models.WebhookEvent) {
    // Your custom logic here
    // Example: Save to database
    // db.SaveMessage(message)
    
    h.logger.Info("Message processed", "message_id", message.ID)
}
```

## 🐳 Deployment

### Production Deployment on Ubuntu VPS

For production deployment on your Ubuntu VPS with domain `asllmarket.org`:

#### 1. Initial VPS Setup

```bash
# On your VPS, run as root:
wget https://raw.githubusercontent.com/your-repo/novinhub-webhook/main/deploy/scripts/setup-vps.sh
chmod +x setup-vps.sh
sudo ./setup-vps.sh
```

This script will:
- Install Go, Nginx, Supervisor, Certbot
- Configure firewall and security
- Set up SSL certificate for `asllmarket.org`
- Create application user and directories
- Configure Nginx reverse proxy

#### 2. Deploy Application

```bash
# On your local machine, upload files to VPS
scp -r . user@asllmarket.org:/tmp/novinhub-webhook/

# SSH into your VPS
ssh user@asllmarket.org

# Navigate to the uploaded directory
cd /tmp/novinhub-webhook

# Run deployment script
chmod +x deploy/scripts/deploy.sh
./deploy/scripts/deploy.sh
```

#### 3. Verify Deployment

```bash
# Check service status
sudo supervisorctl status novinhub-webhook

# Check logs
sudo supervisorctl tail novinhub-webhook

# Test endpoints
curl http://asllmarket.org/health
curl -X POST http://asllmarket.org/webhook -H "Content-Type: application/json" -d '{"type":"test","user_id":"123","payload":{}}'

# Add SSL later (optional)
sudo ./deploy/scripts/add-ssl.sh
```

#### 4. Update Application

```bash
# To update the application
chmod +x deploy/scripts/update.sh
./deploy/scripts/update.sh
```

### Using Docker

```bash
# Build and run with Docker
make docker-build
make docker-run

# Or use Docker Compose
make docker-up
```

### Using Cloud Platforms

The service can be deployed to any cloud platform that supports Go applications:

- **Heroku**: Add a `Procfile` with `web: ./build/webhook`
- **Google Cloud Run**: Deploy as a container
- **AWS Lambda**: Use AWS Lambda Go runtime
- **DigitalOcean App Platform**: Deploy as a web service

## 📊 Monitoring

The service includes:

- Structured JSON logging
- Health check endpoint at `/health`
- Request/response logging
- Error handling and reporting

## 🔒 Security Considerations

- The service returns proper HTTP status codes as required by NovinHub
- CORS is configured for cross-origin requests
- Input validation is performed on all webhook payloads
- Consider adding authentication/authorization for production use

## 🐛 Troubleshooting

### Common Issues

1. **Webhook not receiving events**: Ensure your webhook URL is publicly accessible
2. **Timeout errors**: The service responds within 3 seconds as required by NovinHub
3. **Invalid JSON errors**: Check that NovinHub is sending properly formatted JSON

### Logs

The service logs all webhook events with structured JSON. Check logs for:
- Incoming webhook events
- Processing errors
- Response status codes

## 📝 License

MIT License - feel free to use and modify as needed.
