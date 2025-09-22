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
	SMS      SMSConfig         `mapstructure:"sms"`
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

// SMSConfig holds SMS-related configuration
type SMSConfig struct {
	Provider string        `mapstructure:"provider"`
	IPPanel  IPPanelConfig `mapstructure:"ippanel"`
	Enabled  bool          `mapstructure:"enabled"`
	Retry    RetryConfig   `mapstructure:"retry"`
	Patterns PatternConfig `mapstructure:"patterns"`
}

// IPPanelConfig holds IPPanel-specific configuration
type IPPanelConfig struct {
	APIKey     string `mapstructure:"api_key"`
	Originator string `mapstructure:"originator"`
}

// RetryConfig holds retry-related configuration
type RetryConfig struct {
	MaxAttempts  int `mapstructure:"max_attempts"`
	DelaySeconds int `mapstructure:"delay_seconds"`
}

// PatternConfig holds pattern management configuration
type PatternConfig struct {
	Enabled bool     `mapstructure:"enabled"`
	List    []string `mapstructure:"list"`
	Current int      `mapstructure:"current"`
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
	viper.AddConfigPath("./internal/config")
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

	// SMS defaults
	viper.SetDefault("sms.provider", "ippanel")
	viper.SetDefault("sms.enabled", false)
	viper.SetDefault("sms.retry.max_attempts", 3)
	viper.SetDefault("sms.retry.delay_seconds", 5)

	// Pattern defaults
	viper.SetDefault("sms.patterns.enabled", true)
	viper.SetDefault("sms.patterns.list", []string{
		"a2xjmxbszf27a7e", // گروه اول
		"m3p3jtuu13i4n1o", // گروه دوم
		"l05j64348i04cx8", // گروه سوم
		"nv4fgs9mczuv6rq", // گروه چهارم
	})
	viper.SetDefault("sms.patterns.current", 0)

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

// Pattern Management Methods

// GetCurrentPattern returns the current pattern
func (c *Config) GetCurrentPattern() string {
	if !c.SMS.Patterns.Enabled || len(c.SMS.Patterns.List) == 0 {
		return "" // No patterns configured
	}

	// Always use pattern at current index (0-based)
	if c.SMS.Patterns.Current >= len(c.SMS.Patterns.List) {
		c.SMS.Patterns.Current = 0 // Reset to first pattern if out of bounds
	}

	return c.SMS.Patterns.List[c.SMS.Patterns.Current]
}

// GetCurrentPatternInfo returns current pattern with index and group name
func (c *Config) GetCurrentPatternInfo() (string, int, string) {
	if !c.SMS.Patterns.Enabled || len(c.SMS.Patterns.List) == 0 {
		return "", 0, "هیچ پترنی تنظیم نشده"
	}

	// Ensure current index is within bounds
	if c.SMS.Patterns.Current >= len(c.SMS.Patterns.List) {
		c.SMS.Patterns.Current = 0
	}

	groupNames := []string{"گروه اول", "گروه دوم", "گروه سوم", "گروه چهارم"}
	index := c.SMS.Patterns.Current + 1
	groupName := groupNames[c.SMS.Patterns.Current]

	return c.SMS.Patterns.List[c.SMS.Patterns.Current], index, groupName
}

// NextPattern moves to the next pattern
func (c *Config) NextPattern() (string, int, string) {
	if !c.SMS.Patterns.Enabled || len(c.SMS.Patterns.List) == 0 {
		return "", 0, "هیچ پترنی تنظیم نشده"
	}

	// Move to next pattern (circular)
	c.SMS.Patterns.Current = (c.SMS.Patterns.Current + 1) % len(c.SMS.Patterns.List)

	groupNames := []string{"گروه اول", "گروه دوم", "گروه سوم", "گروه چهارم"}
	index := c.SMS.Patterns.Current + 1
	groupName := groupNames[c.SMS.Patterns.Current]

	return c.SMS.Patterns.List[c.SMS.Patterns.Current], index, groupName
}

// SetPattern sets a specific pattern by index
func (c *Config) SetPattern(index int) error {
	if !c.SMS.Patterns.Enabled || len(c.SMS.Patterns.List) == 0 {
		return fmt.Errorf("pattern management is disabled")
	}

	if index < 0 || index >= len(c.SMS.Patterns.List) {
		return fmt.Errorf("invalid pattern index: %d", index)
	}

	c.SMS.Patterns.Current = index
	return nil
}

// GetPatternsList returns all patterns with their info
func (c *Config) GetPatternsList() []map[string]interface{} {
	if !c.SMS.Patterns.Enabled || len(c.SMS.Patterns.List) == 0 {
		return []map[string]interface{}{
			{
				"index":      0,
				"name":       "هیچ پترنی تنظیم نشده",
				"pattern":    "",
				"is_current": true,
			},
		}
	}

	// Ensure current index is within bounds
	if c.SMS.Patterns.Current >= len(c.SMS.Patterns.List) {
		c.SMS.Patterns.Current = 0
	}

	groupNames := []string{"گروه اول", "گروه دوم", "گروه سوم", "گروه چهارم"}
	patterns := make([]map[string]interface{}, len(c.SMS.Patterns.List))

	for i, pattern := range c.SMS.Patterns.List {
		patterns[i] = map[string]interface{}{
			"index":      i + 1,
			"name":       groupNames[i],
			"pattern":    pattern,
			"is_current": i == c.SMS.Patterns.Current,
		}
	}

	return patterns
}
