package config

import (
	"context"
	"go.uber.org/config"
	"go.uber.org/fx"
)

type Config struct {
	Name string `yaml:"name"`
}

type AppConfig struct {
	fx.Out

	Provider config.Provider
	Config   Config
}

func New(ctx context.Context) (AppConfig, error) {
	cfg := Config{
		Name: "default",
	}
	path := ctx.Value("configPath").(string)
	loader, err := config.NewYAML(config.File(path))
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
