package mongodb

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

func NewMongoDBConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	if err := provider.Get("mongodb").Populate(&cfg); err != nil {
		return nil, fmt.Errorf("mongodb config: %w", err)
	}
	return &cfg, nil
}
