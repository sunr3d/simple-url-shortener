package config

import (
	"fmt"
	"strings"

	"github.com/wb-go/wbf/config"
)

func GetConfig(path string) (*Config, error) {
	cfg := config.New()
	err := cfg.Load(path)
	if err != nil {
		return nil, fmt.Errorf("cfg.Load: %w", err)
	}

	cfg.SetDefault("HTTP_PORT", "8080")
	cfg.SetDefault("BASE_URL", "http://localhost:8080")
	cfg.SetDefault("LOG_LEVEL", "info")

	var c Config
	if err := cfg.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("cfg.Unmarshal: %w", err)
	}
	if strings.TrimSpace(c.DB.DSN) == "" {
		return nil, fmt.Errorf("DB.DSN не может быть пустым")
	}

	return &c, nil
}
