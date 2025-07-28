package config

import (
    "log"
    "github.com/spf13/viper"
)

type AppConfig struct {
    ServerPort string `mapstructure:"SERVER_PORT"`
    TelegramToken string `mapstructure:"TELEGRAM_TOKEN"`
}

func LoadConfig() (*AppConfig, error) {
    viper.SetConfigName("config")
    viper.SetConfigType("yaml")
    viper.AddConfigPath("./internal/config")
    viper.AutomaticEnv()

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    var cfg AppConfig
    if err := viper.Unmarshal(&cfg); err != nil {
        return nil, err
    }

    // Validation
    if cfg.ServerPort == "" {
        log.Fatal("SERVER_PORT must be set in config.yaml or env")
    }
    if cfg.TelegramToken == "" {
        log.Fatal("TELEGRAM_TOKEN must be set in config.yaml or env")
    }

    return &cfg, nil
}