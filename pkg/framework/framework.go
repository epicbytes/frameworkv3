package framework

import (
	"context"
	"errors"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type service struct {
	Name            string
	Version         string
	assignedPlugins map[string]string
	MongoClient     *mongo.Client
	tasks           []runtime.Task
	OnStartup       func(ctx context.Context)
	OnFinish        func(ctx context.Context)
}

func (s service) Run() error {

	log.Info().Msg("Service running")

	if len(s.tasks) == 0 {
		return errors.New("service has no tasks for running")
	}

	var keeper = runtime.TaskKeeper{
		Tasks:           s.tasks,
		ShutdownTimeout: time.Second * 10,
		PingPeriod:      time.Millisecond * 500,
	}

	var app = runtime.Application{
		MainFunc: func(ctx context.Context, halt <-chan struct{}) error {
			var errShutdown = make(chan error, 1)
			if s.OnStartup != nil {
				s.OnStartup(ctx)
			}
			defer func() {
				if s.OnFinish != nil {
					s.OnFinish(ctx)
				}
			}()
			go func() {
				defer close(errShutdown)
				select {
				case <-halt:
				case <-ctx.Done():

				}
			}()
			err, ok := <-errShutdown
			if ok {
				return err
			}
			return nil
		},
		Resources:          &keeper,
		TerminationTimeout: time.Second * 10,
	}
	return app.Run()
}

type Service interface {
	Run() error
}

type ServiceOptions struct {
	Name    string
	Version string
	tasks   []runtime.Task
}

type ServiceOption func(*ServiceOptions)

func NewServiceBuilder(options ...ServiceOption) (srv Service, err error) {
	opts := &ServiceOptions{}
	for _, option := range options {
		option(opts)
	}

	return &service{
		Name:    opts.Name,
		Version: opts.Version,
		tasks:   opts.tasks,
	}, nil

}

func WithNameAndVersion(name, version string) ServiceOption {
	return func(po *ServiceOptions) {
		po.Name = name
		po.Version = version
	}
}

func WithServiceTask(task runtime.Task) ServiceOption {
	return func(po *ServiceOptions) {
		po.tasks = append(po.tasks, task)
	}
}
