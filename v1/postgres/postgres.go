package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/migrate"
	"go.uber.org/zap"
)

type Postgres struct {
	Conn       bun.Conn
	DB         *bun.DB
	Config     *Config
	Logger     *zap.Logger
	Migrations *migrate.Migrations
	Done       chan struct{}
}

func NewPostgres(logger *zap.Logger, config *Config) *Postgres {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.User, config.Password, config.Host, config.Port, config.Database)
	db := bun.NewDB(sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn))), pgdialect.New())
	dbConn, err := db.Conn(context.Background())
	if err != nil {
		logger.Error(err.Error())
	}
	return &Postgres{
		DB:     db,
		Config: config,
		Logger: logger,
		Conn:   dbConn,
		Done:   make(chan struct{}),
	}
}

func (p *Postgres) StartMigrations() error {
	if p.Migrations == nil {
		return nil
	}
	mgrtr := migrate.NewMigrator(p.DB, p.Migrations)
	err := mgrtr.Init(context.Background())
	if err != nil {
		p.Logger.Error(err.Error())
	}
	if err := mgrtr.Lock(context.Background()); err != nil {
		p.Logger.Error(err.Error())
	}
	group, err := mgrtr.Migrate(context.Background())
	if err != nil {
		p.Logger.Warn(err.Error())
		mgrtr.Unlock(context.Background())
		return err
	}
	if group.IsZero() {
		p.Logger.Warn("there are no new migrations to run (database is up to date)")
		mgrtr.Unlock(context.Background())
		return nil
	}
	p.Logger.Info(fmt.Sprintf("migrated to %s\\n", group))
	mgrtr.Unlock(context.Background())
	return nil
}
