package http_server

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Server struct {
	App    *fiber.App
	Config *Config
	Done   chan struct{}
}

func NewServer(config *Config) *Server {

	app := fiber.New(fiber.Config{
		ServerHeader: "EpicServer",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	})

	server := &Server{
		App:    app,
		Config: config,
		Done:   make(chan struct{}),
	}

	return server
}

func (s *Server) StartServer() error {

	s.App.Get("/health", func(ctx *fiber.Ctx) error {
		return ctx.Status(fiber.StatusOK).SendString("ok")
	})
	s.App.Use(otelfiber.Middleware())
	s.App.Use(requestid.New())
	s.App.Use(compress.New(compress.Config{
		Level: compress.LevelBestCompression,
	}))
	s.App.Use(helmet.New(helmet.Config{
		XFrameOptions:             "DENY",
		CrossOriginEmbedderPolicy: "false",
		XSSProtection:             "0",
	}))

	return s.App.Listen(s.Config.Address)
}

func (s *Server) StopServer() error {
	server := s.App
	return server.Shutdown()
}
