// THIS FILE CREATED WITH GENERATOR DO NOT EDIT!
package temporal_worker

import (
	"fmt"
	config "go.uber.org/config"
)

type Config struct {
	TaskQueue string `yaml:"taskQueue"`
}

func NewWorkerConfig(provider config.Provider) (*Config, error) {
	var cfg Config
	if err := provider.Get("worker").Populate(&cfg); err != nil {
		return nil, fmt.Errorf("worker config: %w", err)
	}
	return &cfg, nil
}
