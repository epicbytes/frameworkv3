package postgres

import (
	"context"
	"database/sql"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	pgdialect "github.com/uptrace/bun/dialect/pgdialect"
	pgdriver "github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
)

type OnConnectHandler func(ctx context.Context) error
type Storage interface {
	runtime.Task
}

type storage struct {
	ctx        context.Context
	cfg        *config.Config
	db         *bun.DB
	connection bun.Conn
	migrations *migrate.Migrations
	OnConnect  OnConnectHandler
}

type PostgresClient interface {
	GetClient() bun.Conn
}

func (t *storage) GetClient() bun.Conn {
	return t.connection
}

func New(ctx context.Context, cfg *config.Config, migrations *migrate.Migrations) Storage {
	return &storage{
		ctx:        ctx,
		cfg:        cfg,
		migrations: migrations,
	}
}

func (t *storage) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL Postgres")

	t.db = bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(t.cfg.Postgres.URI))), pgdialect.New())
	t.connection, err = t.db.Conn(context.Background())
	if err != nil {
		log.Warn().Msg(err.Error())
	}

	err = t.runMigrator(t.migrations)
	if err != nil {
		return err
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
	return t.connection.PingContext(context.Background())
}

func (t *storage) Close() error {
	log.Debug().Msg("CLOSE Postgres connection")
	return t.connection.Close()
}
