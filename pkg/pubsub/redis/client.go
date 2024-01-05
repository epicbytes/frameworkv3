package redis

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type redisClient struct {
	ctx       context.Context
	client    *redis.Client
	URI       string
	Password  string
	onConnect func(ctx context.Context, client *redis.Client) error
}

func (t *redisClient) IsPriority() bool {
	return true
}

type RedisClient interface {
	GetClient() *redis.Client
	OnConnect(fn func(ctx context.Context, client *redis.Client) error)
}

func New(cfg *config.Config) runtime.Task {
	return &redisClient{
		URI:      cfg.Redis.URI,
		Password: cfg.Redis.Password,
	}
}

func (t *redisClient) OnConnect(fn func(ctx context.Context, client *redis.Client) error) {
	t.onConnect = fn
}

func (t *redisClient) Init(ctx context.Context) error {
	t.ctx = ctx
	var err error
	log.Debug().Msg("INITIAL Redis")
	t.client = redis.NewClient(&redis.Options{
		Addr:     t.URI,
		Password: t.Password,
		DB:       0,
	})

	if t.onConnect != nil {
		err = t.onConnect(t.ctx, t.client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (t *redisClient) GetClient() *redis.Client {
	return t.client
}

func (t *redisClient) Ping(context.Context) error {
	return nil
}

func (t *redisClient) Close() error {
	log.Debug().Msg("CLOSE Redis connection")
	t.client.Shutdown(t.ctx)
	return nil
}
