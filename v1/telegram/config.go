package telegram

import (
	"fmt"
	"go.uber.org/config"
)

type Config struct {
	Token  string `yaml:"token"`
	Debug  bool   `yaml:"debug"`
	ChatId int64  `yaml:"chatId"`
}

func NewTelegramConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	if err := provider.Get("telegram").Populate(&cfg); err != nil {
		return nil, fmt.Errorf("telegram config: %w", err)
	}
	return &cfg, nil
}
