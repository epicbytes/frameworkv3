package context

import (
	"context"
	"sync"

	"go.uber.org/fx"
)

var once sync.Once

// NewModule fx module for context.
func NewModule() fx.Option {
	options := fx.Options()
	ctx := func() context.Context {
		return context.WithValue(context.Background(), "configPath", "./config.yml")
	}
	once.Do(func() {
		options = fx.Options(
			fx.Provide(
				ctx,
			),
		)
	})

	return options
}
