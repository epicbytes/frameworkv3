package mongodb

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OnConnectHandler func(ctx context.Context) error

type Storage interface {
	runtime.Task
}

type storage struct {
	ctx       context.Context
	URI       string
	DBName    string
	Config    *mgm.Config
	OnConnect OnConnectHandler
}

func New(ctx context.Context, uri string) Storage {
	return &storage{
		ctx: ctx,
		URI: uri,
	}
}

func (t *storage) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL MongoDB")
	err = mgm.SetDefaultConfig(t.Config, t.DBName, options.Client().ApplyURI(t.URI))
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
	return nil //mgm.t.Client.Disconnect(t.ctx)
}
