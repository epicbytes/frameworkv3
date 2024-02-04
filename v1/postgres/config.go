package postgres

import (
	"fmt"
	"go.uber.org/config"
)

type Config struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func NewPostgresConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	if err := provider.Get("postgres").Populate(&cfg); err != nil {
		return nil, fmt.Errorf("postgres config: %w", err)
	}
	return &cfg, nil
}
