package temporal

import (
	"fmt"
	"go.uber.org/config"
)

type Config struct {
	Host      string `yaml:"host"`
	Namespace string `yaml:"namespace"`
}

func NewTemporalConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	if err := provider.Get("temporal").Populate(&cfg); err != nil {
		return nil, fmt.Errorf("temporal config: %w", err)
	}
	return &cfg, nil
}
