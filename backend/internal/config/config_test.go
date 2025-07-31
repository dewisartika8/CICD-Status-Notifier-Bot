package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

const testConfigPath = "../../config/config.yaml"
const testConfigBackup = "../../config/config.yaml.test_backup"

func TestLoadConfigSuccess(t *testing.T) {
	// Backup original config
	if _, err := os.Stat(testConfigPath); err == nil {
		_ = os.Rename(testConfigPath, testConfigBackup)
		defer func() {
			_ = os.Remove(testConfigPath)
			_ = os.Rename(testConfigBackup, testConfigPath)
		}()
	} else {
		defer os.Remove(testConfigPath)
	}

	// Setup: create a temporary config.yaml with new structure
	content := []byte(`server:
  port: 8081
  host: "localhost"
  read_timeout: "30s"
  write_timeout: "30s"
database:
  host: "localhost"
  port: "5432"
  user: "postgres"
  password: "password"
  dbname: "cicd_notifier"
  sslmode: "disable"
  max_open_conns: 10
  max_idle_conns: 5
  max_lifetime: 300
telegram:
  bot_token: "dummy-token"
  webhook_url: "https://example.com/webhook"
github:
  webhook_secret: "secret"
gitlab:
  webhook_secret: "secret"
logging:
  level: "info"
  format: "json"
  output: "stdout"
  file_path: "logs/app.log"
environment: "development"
`)
	_ = os.WriteFile(testConfigPath, content, 0644)

	viper.Reset() // Reset Viper state

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, 8081, cfg.Server.Port)
	assert.Equal(t, "dummy-token", cfg.Telegram.BotToken)
}

func TestLoadConfigMissingServerPort(t *testing.T) {
	// Backup original config
	if _, err := os.Stat(testConfigPath); err == nil {
		_ = os.Rename(testConfigPath, testConfigBackup)
		defer func() {
			_ = os.Remove(testConfigPath)
			_ = os.Rename(testConfigBackup, testConfigPath)
		}()
	} else {
		defer os.Remove(testConfigPath)
	}

	// Setup: create config without server port
	content := []byte(`telegram:
  bot_token: "dummy-token"
`)
	_ = os.WriteFile(testConfigPath, content, 0644)

	viper.Reset()

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, 8080, cfg.Server.Port) // Default port
}

func TestLoadConfigMissingTelegramToken(t *testing.T) {
	// Backup original config
	if _, err := os.Stat(testConfigPath); err == nil {
		_ = os.Rename(testConfigPath, testConfigBackup)
		defer func() {
			_ = os.Remove(testConfigPath)
			_ = os.Rename(testConfigBackup, testConfigPath)
		}()
	} else {
		defer os.Remove(testConfigPath)
	}

	// Setup: create config without telegram token
	content := []byte(`server:
  port: 8081
`)
	_ = os.WriteFile(testConfigPath, content, 0644)

	viper.Reset()

	_, err := LoadConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "config validation failed")
}
