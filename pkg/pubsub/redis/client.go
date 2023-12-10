package redis

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
)

type RedisClient struct {
	ctx       context.Context
	client    *redis.Client
	URI       string
	Password  string
	onConnect func(ctx context.Context, client *redis.Client) error
}

func (t *RedisClient) OnConnect(fn func(ctx context.Context, client *redis.Client) error) {
	t.onConnect = fn
}

func (t *RedisClient) Init(ctx context.Context) error {
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

func (t *RedisClient) GetClient() *redis.Client {
	return t.client
}

func (t *RedisClient) Ping(context.Context) error {
	return nil
}

func (t *RedisClient) Close() error {
	log.Debug().Msg("CLOSE Redis connection")
	t.client.Shutdown(t.ctx)
	return nil
}
