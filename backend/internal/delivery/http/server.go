package http

import (
	"backend/config"
	"backend/internal/domain/health_monitor"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/contrib/fiberzap/v2"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"time"
)

type Server struct {
	inner                *fiber.App
	logger               *zap.SugaredLogger
	validation           *validator.Validate
	config               config.ServerParams
	healthMonitorService health_monitor.Service
}

func NewServer(
	config config.ServerParams,
	healthMonitorService health_monitor.Service,
	logger *zap.Logger,
	validation *validator.Validate,
) *Server {
	fiberConfig := fiber.Config{
		IdleTimeout: time.Duration(config.IdleTimeoutSec) * time.Second,
	}
	server := &Server{
		inner:                fiber.New(fiberConfig),
		config:               config,
		logger:               logger.Sugar(),
		healthMonitorService: healthMonitorService,
		validation:           validation,
	}

	// Middleware
	server.inner.Use(fiberzap.New(fiberzap.Config{
		Logger: logger,
	}))

	// Handlers
	server.inner.Group("/api/v1").
		Get("/container", server.getContainersHandler).
		Get("/container/:container_id", server.getHealthLogHandler).
		Post("/container/health", server.postHealthSnapshotHandler)

	return server
}

func (s *Server) StartAsync() {
	go func() {
		addr := fmt.Sprintf(":%d", s.config.Port)
		if err := s.inner.Listen(addr); err != nil {
			s.logger.Fatal(err)
			panic(err)
		}
	}()
}

func (s *Server) Shutdown() error {
	return s.inner.ShutdownWithTimeout(time.Duration(s.config.ShutdownTimeoutSec) * time.Second)
}
