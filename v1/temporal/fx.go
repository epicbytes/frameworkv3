package temporal

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"telegram",
		fx.Provide(
			NewTemporalConfig,
			NewTemporal,
		),
		fx.Invoke(func(lc fx.Lifecycle, tm *Temporal) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return tm.Start()
				},
				OnStop: func(ctx context.Context) error {
					return tm.Stop()
				},
			})
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("temporal")
		}),
	)
}
