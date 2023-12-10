package interfaces

import (
	"context"
)

type Storage interface {
	Connect(ctx context.Context) error
	Disconnect(ctx context.Context) error
}
