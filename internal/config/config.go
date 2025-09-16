package config

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

// Config holds the application configuration
type Config struct {
	Server   ServerConfig      `mapstructure:"server"`
	Logging  LoggingConfig     `mapstructure:"logging"`
	Webhook  WebhookConfig     `mapstructure:"webhook"`
	Security SecurityConfig    `mapstructure:"security"`
	Health   HealthConfig      `mapstructure:"health"`
	Env      EnvironmentConfig `mapstructure:"environment"`
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port         int    `mapstructure:"port"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
	Host         string `mapstructure:"host"`
}

// LoggingConfig holds logging-related configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level"`
	Format   string `mapstructure:"format"`
	Output   string `mapstructure:"output"`
	FilePath string `mapstructure:"file_path"`
}

// WebhookConfig holds webhook-related configuration
type WebhookConfig struct {
	MaxRequestSize       int64 `mapstructure:"max_request_size"`
	ProcessingTimeout    int   `mapstructure:"processing_timeout"`
	EnableRequestLogging bool  `mapstructure:"enable_request_logging"`
}

// SecurityConfig holds security-related configuration
type SecurityConfig struct {
	EnableCORS     bool     `mapstructure:"enable_cors"`
	AllowedOrigins []string `mapstructure:"allowed_origins"`
	RateLimit      int      `mapstructure:"rate_limit"`
}

// HealthConfig holds health check configuration
type HealthConfig struct {
	Endpoint string `mapstructure:"endpoint"`
	Timeout  int    `mapstructure:"timeout"`
}

// EnvironmentConfig holds environment-specific configuration
type EnvironmentConfig struct {
	Mode  string `mapstructure:"mode"`
	Debug bool   `mapstructure:"debug"`
}

// Load loads configuration from YAML file and environment variables
func Load() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/novinhub-webhook")
	viper.AddConfigPath("$HOME/.novinhub-webhook")

	// Set default values
	setDefaults()

	// Enable reading from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
		// Config file not found, use defaults and environment variables
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.port", 8080)
	viper.SetDefault("server.read_timeout", 10)
	viper.SetDefault("server.write_timeout", 10)
	viper.SetDefault("server.host", "0.0.0.0")

	// Logging defaults
	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.file_path", "/var/log/novinhub-webhook/webhook.log")

	// Webhook defaults
	viper.SetDefault("webhook.max_request_size", 1048576) // 1MB
	viper.SetDefault("webhook.processing_timeout", 30)
	viper.SetDefault("webhook.enable_request_logging", true)

	// Security defaults
	viper.SetDefault("security.enable_cors", true)
	viper.SetDefault("security.allowed_origins", []string{})
	viper.SetDefault("security.rate_limit", 100)

	// Health defaults
	viper.SetDefault("health.endpoint", "/health")
	viper.SetDefault("health.timeout", 5)

	// Environment defaults
	viper.SetDefault("environment.mode", "development")
	viper.SetDefault("environment.debug", false)
}

// GetServerAddress returns the server address in host:port format
func (c *Config) GetServerAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Host, c.Server.Port)
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Env.Mode == "production"
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Env.Mode == "development"
}
