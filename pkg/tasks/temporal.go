package tasks

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/config"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type TemporalClient interface {
	GetClient() client.Client
}

func New(ctx context.Context, cfg *config.Config) runtime.Task {
	return &temporalTask{
		ctx:        ctx,
		URI:        cfg.Temporal.URI,
		Namespaces: cfg.Temporal.Namespaces,
	}
}

type temporalTask struct {
	ctx        context.Context
	URI        string
	Namespaces []string
	Name       string
	Worker     worker.Worker
	client     client.Client
}

func (t *temporalTask) GetClient() client.Client {
	return t.client
}

func (t *temporalTask) Init(context.Context) error {
	log.Debug().Msgf("INITIAL Worker %s", t.Name)
	logger := NewZerologAdapter()
	cl, err := client.Dial(client.Options{
		HostPort:  t.URI,
		Namespace: t.Namespaces[0],
		Logger:    logger,
	})
	t.client = cl
	err = t.Worker.Run(worker.InterruptCh())
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	return nil
}

func (t *temporalTask) Ping(context.Context) error {
	return nil
}

func (t *temporalTask) Close() error {
	log.Debug().Msgf("CLOSE Worker %s", t.Name)
	t.Worker.Stop()
	return nil
}
