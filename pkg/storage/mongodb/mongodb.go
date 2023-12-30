package mongodb

import (
	"context"
	"fmt"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/kamva/mgm/v3"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type OnConnectHandler func(ctx context.Context) error
type BaseModel struct {
	ID        int64 `json:"id" bson:"_id,omitempty"`
	CreatedAt int64 `json:"created_at" bson:"created_at"`
	UpdatedAt int64 `json:"updated_at" bson:"updated_at"`
	DeletedAt int64 `json:"deleted_at" bson:"deleted_at"`
}

func (b *BaseModel) Creating(collName string) error {
	b.CreatedAt = time.Now().Unix()
	b.UpdatedAt = time.Now().Unix()
	var counter struct {
		ID    string `bson:"_id"`
		Value int64  `bson:"value"`
	}
	res := mgm.CollectionByName("counter").FindOneAndUpdate(context.Background(),
		bson.D{{"_id", collName}},
		bson.M{"$inc": bson.M{"value": 1}},
		options.FindOneAndUpdate().SetUpsert(true),
		options.FindOneAndUpdate().SetReturnDocument(options.After))
	if err := res.Err(); err != nil {
		return fmt.Errorf("failed to find one and update: %w", err)
	}
	if err := res.Decode(&counter); err != nil {
		return fmt.Errorf("failed to decode counter: %w", err)
	}
	b.SetID(counter.Value)
	return nil
}

func (b *BaseModel) PrepareID(id interface{}) (interface{}, error) {
	return id, nil
}
func (b *BaseModel) GetID() interface{} {
	return b.ID
}
func (b *BaseModel) SetID(id interface{}) {
	b.ID = id.(int64)
}
func (b *BaseModel) Updating(context.Context) error {
	b.UpdatedAt = time.Now().Unix()
	return nil
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
