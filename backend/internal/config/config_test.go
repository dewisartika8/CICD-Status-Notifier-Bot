package config

import (
	"os"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_Success(t *testing.T) {
	// Setup: create a temporary config.yaml
	content := []byte("SERVER_PORT: \"8081\"\nTELEGRAM_TOKEN: \"dummy-token\"\n")
	_ = os.WriteFile("internal/config/config.yaml", content, 0644)
	defer os.Remove("internal/config/config.yaml")

	viper.Reset() // Reset Viper state

	cfg, err := LoadConfig()
	assert.NoError(t, err)
	assert.Equal(t, "8081", cfg.ServerPort)
	assert.Equal(t, "dummy-token", cfg.TelegramToken)
}

func TestLoadConfig_MissingServerPort(t *testing.T) {
	content := []byte("TELEGRAM_TOKEN: \"dummy-token\"\n")
	_ = os.WriteFile("internal/config/config.yaml", content, 0644)
	defer os.Remove("internal/config/config.yaml")

	viper.Reset()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing SERVER_PORT")
		}
	}()
	_, _ = LoadConfig()
}

func TestLoadConfig_MissingTelegramToken(t *testing.T) {
	content := []byte("SERVER_PORT: \"8081\"\n")
	_ = os.WriteFile("internal/config/config.yaml", content, 0644)
	defer os.Remove("internal/config/config.yaml")

	viper.Reset()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected panic for missing TELEGRAM_TOKEN")
		}
	}()
	_, _ = LoadConfig()
}
