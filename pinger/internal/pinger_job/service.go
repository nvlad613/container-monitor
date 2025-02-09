package pinger_job

import (
	"context"
	"go.uber.org/zap"
	"pinger/config"
	"pinger/internal/api_client"
	"pinger/internal/container_finder"
	"pinger/internal/container_pinger"
	"sync"
	"time"
)

type PingerJob interface {
	Start() error
}

type PingerJobIml struct {
	logger          *zap.SugaredLogger
	client          api_client.ApiClient
	containerFinder container_finder.Finder
	containerPinger container_pinger.Pinger
	conf            config.Config
}

func New(
	logger *zap.SugaredLogger,
	client api_client.ApiClient,
	containerFinder container_finder.Finder,
	containerPinger container_pinger.Pinger,
	conf config.Config,
) *PingerJobIml {
	return &PingerJobIml{
		logger:          logger,
		client:          client,
		containerFinder: containerFinder,
		containerPinger: containerPinger,
		conf:            conf,
	}
}

func (p *PingerJobIml) Start(ctx context.Context) error {
	period := time.Duration(p.conf.Pinger.PingPeriodSec) * time.Second
	for {
		startTime := time.Now()
		if err := p.startWorkers(ctx); err != nil {
			p.logger.Error(err)
		}

		waitTime := calculateWaitTime(startTime, period)
		select {
		case <-ctx.Done():
			return nil
		case <-time.After(waitTime):
		}
	}
}

func (p *PingerJobIml) startWorkers(ctx context.Context) error {
	containers, err := p.containerFinder.ContainersList(ctx)
	if err != nil {
		return err
	}

	workersPool := make(chan struct{}, p.conf.Pinger.Workers)
	resultChan := make(chan container_pinger.HealthReport, p.conf.Pinger.Workers)
	defer close(workersPool)

	// Start workers pool
	go func(pool chan struct{}, results chan<- container_pinger.HealthReport) {
		defer close(results)
		wg := sync.WaitGroup{}

		for _, cnt := range containers {
			wg.Add(1)
			pool <- struct{}{}
			go func() {
				results <- p.containerPinger.Ping(cnt)
				<-pool
				wg.Done()
			}()
		}
		wg.Wait()
	}(workersPool, resultChan)

	for report := range resultChan {
		p.logger.Infof("%s: '%s:%s' %s", report.Name, report.IpAddr, report.Port, report.Status)

		if err = p.client.Post(report, ctx); err != nil {
			p.logger.Errorw("failed to send report", "id", report.DockerId, "name", report.Name)
		}
	}

	return nil
}

func calculateWaitTime(opStarted time.Time, period time.Duration) time.Duration {
	workDuration := time.Now().Sub(opStarted)
	waitTime := period - workDuration
	if waitTime < 0 {
		waitTime = 0
	}

	return waitTime
}
