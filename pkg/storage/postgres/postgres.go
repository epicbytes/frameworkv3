package postgres

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
)

type OnConnectHandler func(ctx context.Context) error
type Storage interface {
	runtime.Task
}

type storage struct {
	ctx        context.Context
	cfg        *config.Config
	connection bun.Conn
	migrations *migrate.Migrations
	OnConnect  OnConnectHandler
}

func (t *storage) IsPriority() bool {
	return true
}

func New(ctx context.Context, connection bun.Conn) Storage {
	return &storage{
		ctx:        ctx,
		connection: connection,
	}
}

type PageEntity struct {
	Id          int64  `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty" bson:"id,omitempty"`
	CreatedAt   int64  `protobuf:"varint,2,opt,name=created_at,json=createdAt,proto3" json:"created_at,omitempty" bson:"created_at,omitempty"`
	UpdatedAt   int64  `protobuf:"varint,3,opt,name=updated_at,json=updatedAt,proto3" json:"updated_at,omitempty" bson:"updated_at,omitempty"`
	DeletedAt   int64  `protobuf:"varint,4,opt,name=deleted_at,json=deletedAt,proto3" json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`
	Name        string `protobuf:"bytes,5,opt,name=name,proto3" json:"name,omitempty" bson:"name,omitempty"`
	Description string `protobuf:"bytes,7,opt,name=description,proto3" json:"description,omitempty" bson:"description,omitempty"`
}

type Page struct {
	bun.BaseModel
	ID        int64      `pathToBun:"id,pk,autoincrement"`
	CreatedAt int64      `pathToBun:",nullzero,notnull,type:bigint"`
	UpdatedAt int64      `pathToBun:",nullzero,notnull,type:bigint"`
	DeletedAt int64      `pathToBun:",soft_delete,nullzero,type:bigint"`
	Entity    PageEntity `bson:"entity" json:"entity"`
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
	return t.connection.PingContext(ctx)
}

func (t *storage) Close() error {
	log.Debug().Msg("CLOSE Postgres connection")
	return t.connection.Close()
}
