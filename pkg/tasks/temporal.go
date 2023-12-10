package tasks

import (
	"context"
	"github.com/rs/zerolog/log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func New(address string, namespace string) (client.Client, error) {
	logger := NewZerologAdapter()
	return client.Dial(client.Options{
		HostPort:  address,
		Namespace: namespace,
		Logger:    logger,
	})
}

type TaskWorker struct {
	Name   string
	Worker worker.Worker
}

func (t *TaskWorker) Init(context.Context) error {
	log.Debug().Msgf("INITIAL Worker %s", t.Name)
	err := t.Worker.Run(worker.InterruptCh())
	if err != nil {
		log.Error().Err(err).Send()
		return err
	}
	return nil
}

func (t *TaskWorker) Ping(context.Context) error {
	return nil
}

func (t *TaskWorker) Close() error {
	log.Debug().Msgf("CLOSE Worker %s", t.Name)
	t.Worker.Stop()
	return nil
}
