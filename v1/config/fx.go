package config

import (
	"go.uber.org/fx"
)

func NewModule() fx.Option {
	return fx.Module(
		"config",
		fx.Provide(
			New,
		),
	)
}
