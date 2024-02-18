package http_server

import (
	"github.com/goccy/go-json"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/contrib/otelfiber"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/healthcheck"
	"github.com/gofiber/fiber/v2/middleware/helmet"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"go.uber.org/zap"
)

type Server struct {
	App    *fiber.App
	Config *Config
	Done   chan struct{}
}

func NewServer(config *Config, handler fiber.ErrorHandler, logger *zap.Logger) *Server {

	cfg := fiber.Config{
		ServerHeader: "EpicServer",
		JSONEncoder:  json.Marshal,
		JSONDecoder:  json.Unmarshal,
	}
	if handler != nil {
		cfg.ErrorHandler = handler
	}
	app := fiber.New(cfg)
	app.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))
	server := &Server{
		App:    app,
		Config: config,
		Done:   make(chan struct{}),
	}

	return server
}

func (s *Server) StartServer() error {

	s.App.Use(otelfiber.Middleware())
	s.App.Use(healthcheck.New(healthcheck.Config{
		LivenessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		LivenessEndpoint: "/live",
		ReadinessProbe: func(c *fiber.Ctx) bool {
			return true
		},
		ReadinessEndpoint: "/ready",
	}))
	s.App.Use(requestid.New())
	s.App.Use(recover.New())
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
