package logger

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"logger",
		fx.Provide(
			NewLogger,
		),
	)
}

func Decorate() fx.Option {
	return fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
		return &fxevent.ZapLogger{Logger: log}
	})
}
