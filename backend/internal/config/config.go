package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/viper"
)

// Constants for environment variables
const (
	EnvServerPort    = "SERVER_PORT"
	EnvTelegramToken = "TELEGRAM_TOKEN"
	EnvDBHost        = "DB_HOST"
	EnvDBPort        = "DB_PORT"
	EnvDBUser        = "DB_USER"
	EnvDBPassword    = "DB_PASSWORD"
	EnvDBName        = "DB_NAME"
)

// Default values
const (
	DefaultEnvironment        = "development"
	DefaultServerPort         = 8080
	DefaultServerHost         = "localhost"
	DefaultServerReadTimeout  = 30 * time.Second
	DefaultServerWriteTimeout = 30 * time.Second

	DefaultDBHost         = "localhost"
	DefaultDBPort         = "5432"
	DefaultDBUser         = "postgres"
	DefaultDBPassword     = "password"
	DefaultDBName         = "cicd_notifier"
	DefaultDBSSLMode      = "disable"
	DefaultDBMaxOpenConns = 10
	DefaultDBMaxIdleConns = 5
	DefaultDBMaxLifetime  = 300 * time.Second

	DefaultLogLevel    = "info"
	DefaultLogFormat   = "json"
	DefaultLogOutput   = "stdout"
	DefaultLogFilePath = "logs/app.log"
)

// ConfigValidationError represents configuration validation errors
type ConfigValidationError struct {
	Field   string
	Message string
}

func (e ConfigValidationError) Error() string {
	return fmt.Sprintf("config validation error for field '%s': %s", e.Field, e.Message)
}

// ConfigLoader interface for loading configuration
type ConfigLoader interface {
	Load() (*AppConfig, error)
}

// ViperConfigLoader implements ConfigLoader using Viper
type ViperConfigLoader struct {
	configPaths []string
	configName  string
	configType  string
}

// NewViperConfigLoader creates a new ViperConfigLoader
func NewViperConfigLoader() *ViperConfigLoader {
	return &ViperConfigLoader{
		configPaths: []string{"./internal/config", "."},
		configName:  "config",
		configType:  "yaml",
	}
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string        `mapstructure:"host" yaml:"host"`
	Port         string        `mapstructure:"port" yaml:"port"`
	User         string        `mapstructure:"user" yaml:"user"`
	Password     string        `mapstructure:"password" yaml:"password"`
	DBName       string        `mapstructure:"dbname" yaml:"dbname"`
	SSLMode      string        `mapstructure:"sslmode" yaml:"sslmode"`
	MaxOpenConns int           `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int           `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxLifetime  time.Duration `mapstructure:"max_lifetime" yaml:"max_lifetime"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         int           `mapstructure:"port" yaml:"port"`
	Host         string        `mapstructure:"host" yaml:"host"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" yaml:"write_timeout"`
}

// TelegramConfig holds telegram bot configuration
type TelegramConfig struct {
	BotToken   string `mapstructure:"bot_token" yaml:"bot_token"`
	WebhookURL string `mapstructure:"webhook_url" yaml:"webhook_url"`
}

// GitHubConfig holds GitHub webhook configuration
type GitHubConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret" yaml:"webhook_secret"`
}

// GitLabConfig holds GitLab webhook configuration
type GitLabConfig struct {
	WebhookSecret string `mapstructure:"webhook_secret" yaml:"webhook_secret"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level    string `mapstructure:"level" yaml:"level"`
	Format   string `mapstructure:"format" yaml:"format"`
	Output   string `mapstructure:"output" yaml:"output"`
	FilePath string `mapstructure:"file_path" yaml:"file_path"`
}

// AppConfig holds all application configuration
type AppConfig struct {
	Environment string         `mapstructure:"environment" yaml:"environment"`
	Server      ServerConfig   `mapstructure:"server" yaml:"server"`
	Database    DatabaseConfig `mapstructure:"database" yaml:"database"`
	Telegram    TelegramConfig `mapstructure:"telegram" yaml:"telegram"`
	GitHub      GitHubConfig   `mapstructure:"github" yaml:"github"`
	GitLab      GitLabConfig   `mapstructure:"gitlab" yaml:"gitlab"`
	Logging     LoggingConfig  `mapstructure:"logging" yaml:"logging"`

	// Backward compatibility
	ServerPort    int    `mapstructure:"SERVER_PORT"`
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
}

// Load implements ConfigLoader interface
func (loader *ViperConfigLoader) Load() (*AppConfig, error) {
	v := viper.New()

	if err := loader.setupViper(v); err != nil {
		return nil, fmt.Errorf("failed to setup viper: %w", err)
	}

	var cfg AppConfig
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	loader.handleBackwardCompatibility(&cfg)
	loader.overrideWithEnvVars(&cfg)

	if err := validateConfig(&cfg); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return &cfg, nil
}

// setupViper configures viper with config paths, defaults, and file reading
func (loader *ViperConfigLoader) setupViper(v *viper.Viper) error {
	v.SetConfigName(loader.configName)
	v.SetConfigType(loader.configType)

	for _, path := range loader.configPaths {
		v.AddConfigPath(path)
	}

	v.AutomaticEnv()
	setDefaults(v)

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	return nil
}

// handleBackwardCompatibility maintains compatibility with old config fields
func (loader *ViperConfigLoader) handleBackwardCompatibility(cfg *AppConfig) {
	if cfg.ServerPort != 0 && cfg.Server.Port == 0 {
		cfg.Server.Port = cfg.ServerPort
	}
	if cfg.TelegramToken != "" && cfg.Telegram.BotToken == "" {
		cfg.Telegram.BotToken = cfg.TelegramToken
	}
}

// overrideWithEnvVars overrides config values with environment variables
func (loader *ViperConfigLoader) overrideWithEnvVars(cfg *AppConfig) {
	// Override string fields
	envOverrides := map[string]*string{
		EnvTelegramToken: &cfg.Telegram.BotToken,
		EnvDBHost:        &cfg.Database.Host,
		EnvDBPort:        &cfg.Database.Port,
		EnvDBUser:        &cfg.Database.User,
		EnvDBPassword:    &cfg.Database.Password,
		EnvDBName:        &cfg.Database.DBName,
	}

	for envKey, configField := range envOverrides {
		if value := os.Getenv(envKey); value != "" {
			*configField = value
		}
	}

	// Override integer fields
	if value := os.Getenv(EnvServerPort); value != "" {
		if port, err := strconv.Atoi(value); err == nil {
			cfg.Server.Port = port
		}
	}
}

// LoadConfig loads configuration using the default ViperConfigLoader
func LoadConfig() (*AppConfig, error) {
	loader := NewViperConfigLoader()
	return loader.Load()
}

// setDefaults sets default configuration values
func setDefaults(v *viper.Viper) {
	v.SetDefault("environment", DefaultEnvironment)
	v.SetDefault("server.port", DefaultServerPort)
	v.SetDefault("server.host", DefaultServerHost)
	v.SetDefault("server.read_timeout", DefaultServerReadTimeout)
	v.SetDefault("server.write_timeout", DefaultServerWriteTimeout)

	v.SetDefault("database.host", DefaultDBHost)
	v.SetDefault("database.port", DefaultDBPort)
	v.SetDefault("database.user", DefaultDBUser)
	v.SetDefault("database.password", DefaultDBPassword)
	v.SetDefault("database.dbname", DefaultDBName)
	v.SetDefault("database.sslmode", DefaultDBSSLMode)
	v.SetDefault("database.max_open_conns", DefaultDBMaxOpenConns)
	v.SetDefault("database.max_idle_conns", DefaultDBMaxIdleConns)
	v.SetDefault("database.max_lifetime", DefaultDBMaxLifetime)

	v.SetDefault("logging.level", DefaultLogLevel)
	v.SetDefault("logging.format", DefaultLogFormat)
	v.SetDefault("logging.output", DefaultLogOutput)
	v.SetDefault("logging.file_path", DefaultLogFilePath)
}

// validateConfig validates the configuration
func validateConfig(cfg *AppConfig) error {
	var validationErrors []error

	// Validate server configuration
	if err := validateServerConfig(&cfg.Server); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate database configuration
	if err := validateDatabaseConfig(&cfg.Database); err != nil {
		validationErrors = append(validationErrors, err)
	}

	// Validate logging configuration
	if err := validateLoggingConfig(&cfg.Logging); err != nil {
		validationErrors = append(validationErrors, err)
	}

	if len(validationErrors) > 0 {
		return combineErrors(validationErrors)
	}

	return nil
}

// validateServerConfig validates server configuration
func validateServerConfig(cfg *ServerConfig) error {
	if cfg.Port == 0 {
		return ConfigValidationError{
			Field:   "server.port",
			Message: "port must be set",
		}
	}

	if cfg.Host == "" {
		return ConfigValidationError{
			Field:   "server.host",
			Message: "host must be set",
		}
	}

	return nil
}

// validateDatabaseConfig validates database configuration
func validateDatabaseConfig(cfg *DatabaseConfig) error {
	required := map[string]string{
		"database.host":   cfg.Host,
		"database.port":   cfg.Port,
		"database.user":   cfg.User,
		"database.dbname": cfg.DBName,
	}

	for field, value := range required {
		if value == "" {
			return ConfigValidationError{
				Field:   field,
				Message: "field is required",
			}
		}
	}

	if cfg.MaxOpenConns <= 0 {
		return ConfigValidationError{
			Field:   "database.max_open_conns",
			Message: "must be greater than 0",
		}
	}

	if cfg.MaxIdleConns < 0 {
		return ConfigValidationError{
			Field:   "database.max_idle_conns",
			Message: "must be greater than or equal to 0",
		}
	}

	return nil
}

// validateLoggingConfig validates logging configuration
func validateLoggingConfig(cfg *LoggingConfig) error {
	validLevels := []string{"debug", "info", "warn", "error", "fatal", "panic"}
	if !contains(validLevels, strings.ToLower(cfg.Level)) {
		return ConfigValidationError{
			Field:   "logging.level",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validLevels, ", ")),
		}
	}

	validFormats := []string{"json", "text"}
	if !contains(validFormats, strings.ToLower(cfg.Format)) {
		return ConfigValidationError{
			Field:   "logging.format",
			Message: fmt.Sprintf("must be one of: %s", strings.Join(validFormats, ", ")),
		}
	}

	return nil
}

// combineErrors combines multiple errors into a single error
func combineErrors(errors []error) error {
	var messages []string
	for _, err := range errors {
		messages = append(messages, err.Error())
	}
	return fmt.Errorf("validation errors: %s", strings.Join(messages, "; "))
}

// contains checks if a slice contains a specific string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
