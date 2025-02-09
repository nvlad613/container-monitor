package main

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"net/http"
	"os"
	"os/signal"
	"pinger/config"
	"pinger/internal/api_client"
	"pinger/internal/container_finder"
	"pinger/internal/container_pinger"
	"pinger/internal/pinger_job"
	"sync"
	"syscall"
)

func main() {
	// Read config
	conf := lo.Must(config.Load())
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}

	// Init logger
	logger := lo.Must(conf.Logger.Build())
	sugaredLogger := logger.Sugar()
	defer logger.Sync()

	// Init validator
	validation := validator.New()

	// Init clients
	dockerClient := lo.Must(client.NewClientWithOpts(
		client.WithHost(conf.Docker.Host),
		client.WithAPIVersionNegotiation(),
	))
	httpClient := http.Client{}
	apiClient := api_client.New(httpClient, conf.ConsumerServer)

	finder := container_finder.New(sugaredLogger, dockerClient, validation)
	pinger := container_pinger.New(conf.Pinger)

	pingerJob := pinger_job.New(
		sugaredLogger,
		apiClient,
		finder,
		pinger,
		conf,
	)

	// Run pinger in background
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := pingerJob.Start(ctx); err != nil {
			sugaredLogger.Fatal(err)
		}
	}()

	// Wait for interrupt
	WaitInterrupt()
	cancel()
	wg.Wait()

	sugaredLogger.Info("Service shutdown")
}

func WaitInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}
