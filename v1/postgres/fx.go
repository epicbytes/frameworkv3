package postgres

import (
	"context"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"log"
)

func NewModule() fx.Option {
	return fx.Module(
		"postgres",
		fx.Provide(
			NewPostgresConfig,
			NewPostgres,
		),
		fx.Invoke(func(lc fx.Lifecycle, pg *Postgres) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					go func() {
						if err := pg.StartMigrations(); err != nil {
							log.Fatalf("start server error : %v\n", err)
						}
					}()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return pg.Conn.Close()
				},
			})
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("postgres")
		}),
	)
}
