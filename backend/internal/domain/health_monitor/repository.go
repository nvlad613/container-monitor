package health_monitor

import (
	"context"
	"time"
)

type Repository interface {
	GetContainersRepository
	InsertHealthSnapshotsRepository
}

type GetContainersRepository interface {
	GetContainers(offset int, limit int, ctx context.Context) ([]Container, error)
	GetHealthSnapshots(containerId int, from time.Time, to time.Time, ctx context.Context) ([]HealthSnapshot, error)
}

type InsertHealthSnapshotsRepository interface {
	InsertSnapshotAndUpdateContainer(report HealthReport, ctx context.Context) error
	InsertActiveSnapshotAndUpdateContainer(report HealthReport, ctx context.Context) error
}
