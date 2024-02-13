package mongodb

import (
	"context"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
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
		fx.Invoke(func(lc fx.Lifecycle, config *Config) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					return mgm.SetDefaultConfig(nil, config.Database, options.Client().ApplyURI("mongodb://"+config.User+":"+config.Password+"@"+config.Host+"/"+config.Database+"?retryWrites=true&replicaSet=dbrs&readPreference=primary&connectTimeoutMS=10000&authSource="+config.Database+"&authMechanism=SCRAM-SHA-1"))
				},
			})
		}),
		fx.Decorate(func(log *zap.Logger) *zap.Logger {
			return log.Named("mongodb")
		}),
	)
}
