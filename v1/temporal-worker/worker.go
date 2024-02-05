// THIS FILE CREATED WITH GENERATOR DO NOT EDIT!
package temporal_worker

import (
	temporal "github.com/epicbytes/frameworkv3/v1/temporal"
	client "go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	zap "go.uber.org/zap"
)

type TemporalWorker struct {
	Config   *Config
	Logger   *zap.Logger
	Temporal client.Client
	Worker   worker.Worker
}

func NewWorker(config *Config, logger *zap.Logger, temporal *temporal.Temporal) *TemporalWorker {
	wrk := &TemporalWorker{
		Config:   config,
		Logger:   logger,
		Temporal: temporal.Client,
		Worker:   worker.New(temporal.Client, config.TaskQueue, worker.Options{}),
	}
	return wrk
}
func (s *TemporalWorker) StartWorker() error {
	return s.Worker.Run(worker.InterruptCh())
}
func (s *TemporalWorker) StopWorker() error {
	s.Worker.Stop()
	return nil
}
