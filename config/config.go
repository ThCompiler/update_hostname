package config

import (
	"fmt"
	"update_hostname/internal/logger"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		Login       string     `yaml:"login"`
		Password    string     `yaml:"password"`
		Domain      string     `yaml:"domain"`
		UpdateHours int64      `yaml:"update_hours"`
		LoggerInfo  LoggerInfo `yaml:"logger"`
	}

	LoggerInfo struct {
		AppName           string          `yaml:"app_name"`
		Directory         string          `yaml:"directory"`
		Level             logger.LogLevel `yaml:"level"`
		UseStdAndFile     bool            `yaml:"use_std_and_file"`
		AllowShowLowLevel bool            `yaml:"allow_show_low_level"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
