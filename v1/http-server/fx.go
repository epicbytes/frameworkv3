package http_server

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

func NewModule() fx.Option {
	return fx.Module(
		"http-server",
		fx.Provide(
			NewServerConfig,
			NewServer,
		),
		fx.Invoke(func(lc fx.Lifecycle, server *Server) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := server.StartServer(); err != nil {
							log.Fatalf("start server error : %v\n", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					go func() {
						err := server.StopServer()
						if err != nil {
							log.Fatalf("stoping server error : %v\n", err)
						}
					}()
					return nil
				},
			})
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("http-server")
		}),
	)
}
