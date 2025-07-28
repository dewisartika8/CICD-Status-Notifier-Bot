package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host         string `mapstructure:"host" yaml:"host"`
	Port         string `mapstructure:"port" yaml:"port"`
	User         string `mapstructure:"user" yaml:"user"`
	Password     string `mapstructure:"password" yaml:"password"`
	DBName       string `mapstructure:"dbname" yaml:"dbname"`
	SSLMode      string `mapstructure:"sslmode" yaml:"sslmode"`
	MaxOpenConns int    `mapstructure:"max_open_conns" yaml:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns" yaml:"max_idle_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime" yaml:"max_lifetime"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string `mapstructure:"port" yaml:"port"`
	Host         string `mapstructure:"host" yaml:"host"`
	ReadTimeout  string `mapstructure:"read_timeout" yaml:"read_timeout"`
	WriteTimeout string `mapstructure:"write_timeout" yaml:"write_timeout"`
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
	ServerPort    string `mapstructure:"SERVER_PORT"`
	TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
}

// LoadConfig loads configuration from file and environment variables
func LoadConfig() (*AppConfig, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./internal/config")
	viper.AddConfigPath(".")

	// Enable environment variable binding
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		// If config file not found, continue with defaults and env vars
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var cfg AppConfig
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	// Handle backward compatibility
	if cfg.ServerPort != "" && cfg.Server.Port == "" {
		cfg.Server.Port = cfg.ServerPort
	}
	if cfg.TelegramToken != "" && cfg.Telegram.BotToken == "" {
		cfg.Telegram.BotToken = cfg.TelegramToken
	}

	// Override with environment variables if present
	if port := os.Getenv("SERVER_PORT"); port != "" {
		cfg.Server.Port = port
	}
	if token := os.Getenv("TELEGRAM_TOKEN"); token != "" {
		cfg.Telegram.BotToken = token
	}
	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.Database.Host = dbHost
	}
	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.Database.Port = dbPort
	}
	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.Database.User = dbUser
	}
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		cfg.Database.Password = dbPassword
	}
	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.Database.DBName = dbName
	}

	// Validate required fields
	if err := validateConfig(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// setDefaults sets default configuration values
func setDefaults() {
	viper.SetDefault("environment", "development")
	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.host", "localhost")
	viper.SetDefault("server.read_timeout", "30s")
	viper.SetDefault("server.write_timeout", "30s")

	viper.SetDefault("database.host", "localhost")
	viper.SetDefault("database.port", "5432")
	viper.SetDefault("database.user", "postgres")
	viper.SetDefault("database.password", "password")
	viper.SetDefault("database.dbname", "cicd_notifier")
	viper.SetDefault("database.sslmode", "disable")
	viper.SetDefault("database.max_open_conns", 10)
	viper.SetDefault("database.max_idle_conns", 5)
	viper.SetDefault("database.max_lifetime", 300)

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "json")
	viper.SetDefault("logging.output", "stdout")
	viper.SetDefault("logging.file_path", "logs/app.log")
}

// validateConfig validates the configuration
func validateConfig(cfg *AppConfig) error {
	if cfg.Server.Port == "" {
		log.Fatal("Server port must be set")
	}

	// Optional validation for other fields
	// Telegram token is optional during development

	return nil
}
