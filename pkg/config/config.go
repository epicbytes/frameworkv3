package config

import (
	"github.com/caarlos0/env/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func New(opts ...Option) (config *Config, err error) {
	cfg := &Config{}

	for _, o := range opts {
		o.apply(cfg)
	}

	if len(cfg.EnvFile.Path) > 0 {
		log.Debug().Msg("loading environment from file")
		err = godotenv.Load(cfg.EnvFile.Path)
		if err != nil {
			log.Warn().Msg("error loading from file, now load from environment")
			//return nil, err
		}
	}

	if err := env.Parse(cfg); err != nil {
		log.Error().Msgf("%+v\n", err)
	}

	return cfg, nil
}
