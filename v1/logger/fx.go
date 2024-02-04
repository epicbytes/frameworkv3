package logger

import "go.uber.org/fx"

func NewModule() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(
			NewLogger,
		),
	)
}
