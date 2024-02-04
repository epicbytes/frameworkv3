package temporal

import (
	"context"
	"go.temporal.io/sdk/client"
	"go.uber.org/zap"
)

type Temporal struct {
	Client client.Client
	Logger *zap.Logger
}

func NewTemporal(ctx context.Context, config *Config, logger *zap.Logger) *Temporal {
	cl, err := client.Dial(client.Options{
		HostPort:  config.Host,
		Namespace: config.Namespace,
		Logger:    NewZapAdapter(logger),
	})
	if err != nil {
		logger.Error(err.Error())
	}

	return &Temporal{
		Logger: logger,
		Client: cl,
	}
}

func (c *Temporal) Start() error {
	return nil
}
func (c *Temporal) Stop() error {
	c.Client.Close()
	return nil
}
