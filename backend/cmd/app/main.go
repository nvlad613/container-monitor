package main

import (
	"backend/config"
	"backend/internal/delivery/http"
	"backend/internal/domain/health_monitor"
	"backend/internal/infra"
	"backend/internal/infra/repository"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"os"
	"os/signal"
	"syscall"
)

// @title           Swagger API
// @version         1.0
// @description     REST API for health monitor

// @host      localhost:3000
// @BasePath  /api/v1
// @schemes http
func main() {
	// Init configuration
	conf := lo.Must(config.Load())

	// Init logger
	logger := lo.Must(conf.Logger.Build())
	sugaredLogger := logger.Sugar()
	defer logger.Sync()

	// Init package structures
	db := lo.Must(infra.InitDb(conf.Database))
	defer db.Close()
	validation := validator.New(validator.WithRequiredStructEnabled())

	// Init repository
	healthMonitorRepository := repository.NewHealthMonitor(db)
	// Init service
	healthMonitorService := health_monitor.New(healthMonitorRepository, sugaredLogger)

	// Init server application
	server := http.NewServer(conf.Server, healthMonitorService, logger, validation)
	server.StartAsync()

	logger.Info("Service started!")
	WaitInterrupt()

	logger.Info("Gracefully shutting down...")
	if err := server.Shutdown(); err != nil {
		sugaredLogger.Errorf("failed to shutdown gracefully: %v", err)
	}

	logger.Info("Server shutdown")
}

func WaitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
