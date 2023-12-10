package postgres

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
)

type OnConnectHandler func(ctx context.Context) error
type Storage interface {
	runtime.Task
}

type storage struct {
	ctx        context.Context
	URI        string
	DBName     string
	connection bun.Conn
	OnConnect  OnConnectHandler
}

func New(ctx context.Context, connection bun.Conn) Storage {
	return &storage{
		ctx:        ctx,
		connection: connection,
	}
}

func (t *storage) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL Postgres")

	if t.OnConnect != nil {
		err = t.OnConnect(t.ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *storage) Ping(ctx context.Context) error {
	return t.connection.PingContext(context.Background())
}

func (t *storage) Close() error {
	log.Debug().Msg("CLOSE Postgres connection")
	return t.connection.Close()
}
