package network

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"net/http"
)

type grpcHTTP struct {
	ctx     context.Context
	handler http.Handler
	path    string
	address string
}

type GrpcHTTP interface {
	runtime.Task
}

func NewGRPCServer(path string, handler http.Handler, address string) GrpcHTTP {
	if address == "" {
		address = ":8080"
	}
	return &grpcHTTP{
		address: address,
		path:    path,
		handler: handler,
	}
}

func (g grpcHTTP) Init(ctx context.Context) error {
	log.Debug().Msg("INITIAL GRPC SERVICE")
	mux := http.NewServeMux()
	mux.Handle(g.path, g.handler)
	log.Info().Msgf("GRPC server started at %s %s", g.path, g.address)
	go func() {
		err := http.ListenAndServe(
			g.address,
			h2c.NewHandler(mux, &http2.Server{}),
		)
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
	}()
	return nil
}

func (g grpcHTTP) Ping(ctx context.Context) error {
	return nil
}

func (g grpcHTTP) Close() error {
	return nil
}
