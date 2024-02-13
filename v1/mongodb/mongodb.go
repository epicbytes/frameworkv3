package mongodb

import (
	"github.com/kamva/mgm/v3"
	"github.com/uptrace/bun/migrate"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

type Mongodb struct {
	Coll       func(m mgm.Model, opts ...*options.CollectionOptions) *mgm.Collection
	Config     *Config
	Logger     *zap.Logger
	Migrations *migrate.Migrations
}

func NewMongoDB(logger *zap.Logger, config *Config) *Mongodb {
	logger.Debug("mongodb://" + config.User + ":" + config.Password + "@" + config.Host + "/" + config.Database + "?retryWrites=true&replicaSet=dbrs&readPreference=primary&connectTimeoutMS=10000&authSource=" + config.Database + "&authMechanism=SCRAM-SHA-1")
	err := mgm.SetDefaultConfig(nil, config.Database, options.Client().ApplyURI("mongodb://"+config.User+":"+config.Password+"@"+config.Host+"/"+config.Database+"?retryWrites=true&replicaSet=dbrs&readPreference=primary&connectTimeoutMS=10000&authSource="+config.Database+"&authMechanism=SCRAM-SHA-1"))
	if err != nil {
		logger.Error(err.Error())
	}
	return &Mongodb{
		Config: config,
		Logger: logger,
		Coll:   mgm.Coll,
	}

}
