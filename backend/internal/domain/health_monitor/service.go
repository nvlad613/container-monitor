package health_monitor

import (
	"context"
	"go.uber.org/zap"
	"time"
)

type Service interface {
	GetContainersList(offset int, limit int, ctx context.Context) ([]Container, error)
	GetContainerHealth(containerId int, from time.Time, to time.Time, ctx context.Context) (*HealthLog, error)
	RecordHealthReport(report HealthReport, ctx context.Context) error
}

type ServiceImpl struct {
	repository Repository
	logger     *zap.SugaredLogger
}

func New(repository Repository, logger *zap.SugaredLogger) *ServiceImpl {
	return &ServiceImpl{
		repository: repository,
		logger:     logger,
	}
}
