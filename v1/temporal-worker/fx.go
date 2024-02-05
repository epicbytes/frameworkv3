// THIS FILE CREATED WITH GENERATOR DO NOT EDIT!
package temporal_worker

import (
	"context"
	fx "go.uber.org/fx"
	zap "go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module("temporal-worker", fx.Provide(NewWorkerConfig, NewWorker), fx.Invoke(func(lc fx.Lifecycle, wr *TemporalWorker) {
		lc.Append(fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					err := wr.StartWorker()
					if err != nil {
						return
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				return wr.StopWorker()
			},
		})
	}), fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("temporal-worker")
	}))
}
