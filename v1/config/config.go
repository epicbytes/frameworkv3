package config

import (
	"context"
	"go.uber.org/config"
	"go.uber.org/fx"
	"io"
)

type Config struct {
	Name string `yaml:"name"`
}

type AppConfig struct {
	fx.Out

	Provider config.Provider
	Config   Config
}

// New creates config for service, is configReader is not null then config will be parsed from config.yml file
func New(ctx context.Context, configReader io.Reader) (AppConfig, error) {
	cfg := Config{
		Name: "default",
	}
	var source config.YAMLOption
	if configReader != nil {
		source = config.Source(configReader)
	} else {
		source = config.Name("./config.yaml")
	}
	loader, err := config.NewYAML(source)
	if err != nil {
		return AppConfig{}, err
	}

	if err := loader.Get("app").Populate(&cfg); err != nil {
		return AppConfig{}, err
	}

	return AppConfig{
		Provider: loader,
		Config:   cfg,
	}, nil
}
