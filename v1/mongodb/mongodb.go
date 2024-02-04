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

	return &Mongodb{
		Config: config,
		Logger: logger,
		Coll:   mgm.Coll,
	}

}
