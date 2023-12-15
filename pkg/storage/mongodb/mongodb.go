package mongodb

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type OnConnectHandler func(ctx context.Context) error
type BaseModel struct {
	mgm.DefaultModel `bson:",inline"`
	DeletedAt        time.Time `bson:"deleted_at,omitempty"`
}

func (b *BaseModel) Deleting(ctx context.Context) error {

	return nil
}

type Storage interface {
	runtime.Task
}
type storage struct {
	ctx          context.Context
	uri          string
	databaseName string
	config       *mgm.Config
	OnConnect    OnConnectHandler
}

func New(ctx context.Context, cfg *config.Config, config *mgm.Config, connectHandler OnConnectHandler) Storage {
	st := &storage{
		ctx:       ctx,
		config:    config,
		OnConnect: connectHandler,
	}

	if cfg != nil && (cfg.Mongo.URI != "" || cfg.Mongo.DatabaseName != "") {
		st.uri = cfg.Mongo.URI
		st.databaseName = cfg.Mongo.DatabaseName
	}

	return st
}

func (t *storage) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL MongoDB")
	if t.uri == "" || t.databaseName == "" {
		log.Fatal().Msg("mongo client config is not set")
		return errors.New("mongo client config is not set")
	}
	err = mgm.SetDefaultConfig(t.config, t.databaseName, options.Client().ApplyURI(t.uri))
	if err != nil {
		return errors.New("mongo client is not set")
	}
	if t.OnConnect != nil {
		err = t.OnConnect(t.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *storage) Ping(ctx context.Context) error {
	return nil //mgm.Ping(ctx, nil)
}

func (t *storage) Close() error {
	log.Debug().Msg("CLOSE MongoDB connection")
	return nil
}
