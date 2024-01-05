package network

import (
	"context"
	"github.com/epicbytes/frameworkv3/pkg/runtime"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

type fiberHTTP struct {
	ctx     context.Context
	app     *fiber.App
	address string
}

func (t *fiberHTTP) IsPriority() bool {
	return false
}

type FiberHTTP interface {
	runtime.Task
}

func NewFiber(app *fiber.App, address string) FiberHTTP {
	if address == "" {
		address = ":8099"
	}
	return &fiberHTTP{
		address: address,
		app:     app,
	}
}

func (t *fiberHTTP) Init(ctx context.Context) error {
	t.ctx = ctx
	log.Debug().Msg("INITIAL HTTP SERVICE")
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	t.app.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))
	t.app.Use(requestid.New(requestid.Config{}))
	//t.server.Use(helmet.New(helmet.Config{}))
	t.app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodOptions,
			fiber.MethodPut,
			fiber.MethodDelete,
			fiber.MethodPatch,
		}, ","),
	}))
	log.Info().Msgf("HTTP server started at %s", t.address)
	go func() {
		err := t.app.Listen(t.address)
		if err != nil {
			log.Error().Err(err).Send()
			return
		}
	}()
	return nil
}

func (t *fiberHTTP) Ping(context.Context) error {
	//log.Debug().Msg("PING GRPC")
	return nil
}

func (t *fiberHTTP) Close() error {
	log.Debug().Msg("CLOSE HTTP")
	err := t.app.Shutdown()
	if err != nil {
		return err
	}
	return nil
}
