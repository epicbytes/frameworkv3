package http

import (
	"context"
	"fmt"
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
	"time"
)

type HTTPService struct {
	ctx         context.Context
	requestTime time.Time
	server      *fiber.App
}

func (t *HTTPService) SetServer(server *fiber.App) {
	t.server = server
}

func (t *HTTPService) GetServer() *fiber.App {
	return t.server
}

func (t *HTTPService) Init(ctx context.Context) error {
	t.ctx = ctx
	log.Debug().Msg("INITIAL HTTP SERVICE")
	logger := zerolog.New(os.Stderr).With().Timestamp().Logger()

	t.server.Use(fiberzerolog.New(fiberzerolog.Config{
		Logger: &logger,
	}))

	t.server.Use(cors.New(cors.Config{
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
	log.Info().Msgf("HTTP server started at %s", ":8099")
	go func() {
		err := t.server.Listen(":8099")
		if err != nil {
			fmt.Print("ERR", err)
			log.Error().Err(err)
			return
		}
	}()
	return nil
}

func (t *HTTPService) Ping(context.Context) error {
	//log.Debug().Msg("PING GRPC")
	return nil
}

func (t *HTTPService) Close() error {
	log.Debug().Msg("CLOSE HTTP")
	err := t.server.Shutdown()
	if err != nil {
		return err
	}
	return nil
}
