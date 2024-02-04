package mqtt

import (
	"context"
	"fmt"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"mqtt",
		fx.Provide(
			NewMqttConfig,
			NewMqtt,
		),
		fx.Invoke(func(lc fx.Lifecycle, mq *MQTT) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := mq.StartMqtt(); err != nil {
							mq.Logger.Fatal(fmt.Sprintf("start server error : %v\n", err))
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return mq.StopMqtt()
				},
			})
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("mqtt")
		}),
	)
}
