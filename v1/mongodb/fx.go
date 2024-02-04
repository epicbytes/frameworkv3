package mongodb

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func NewModule() fx.Option {
	return fx.Module(
		"mongodb",
		fx.Provide(
			NewMongoDBConfig,
			NewMongoDB,
		),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("mongodb")
		}),
	)
}
