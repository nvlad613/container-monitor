package repository

import (
	"backend/internal/domain/health_monitor"
	"context"
	"github.com/samber/lo"
	"github.com/uptrace/bun"
	"time"
)

type HealthMonitorRepositoryImpl struct {
	db *bun.DB
}

func NewHealthMonitor(db *bun.DB) *HealthMonitorRepositoryImpl {
	return &HealthMonitorRepositoryImpl{
		db,
	}
}

func (r *HealthMonitorRepositoryImpl) GetContainers(offset int, limit int, ctx context.Context) ([]health_monitor.Container, error) {
	var containers []Container
	err := r.db.NewSelect().
		Model(&containers).
		Offset(offset).
		Limit(limit).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return lo.Map(containers, func(item Container, _ int) health_monitor.Container {
		return item.ToModel()
	}), nil
}

func (r *HealthMonitorRepositoryImpl) GetHealthSnapshots(containerId int, from time.Time, to time.Time, ctx context.Context) ([]health_monitor.HealthSnapshot, error) {
	var log []HealthSnapshot
	err := r.db.NewSelect().
		Model(&log).
		Where("container_id = ?", containerId).
		Where("timestamp > ?", from).
		Where("timestamp < ?", to).
		Scan(ctx)

	if err != nil {
		return nil, err
	}

	return lo.Map(log, func(item HealthSnapshot, _ int) health_monitor.HealthSnapshot {
		return item.ToModel()
	}), nil
}

func (r *HealthMonitorRepositoryImpl) InsertSnapshotAndUpdateContainer(report health_monitor.HealthReport, ctx context.Context) error {
	container := Container{
		DockerId:   report.DockerId,
		IpAddr:     report.IpAddr,
		Name:       report.Name,
		Status:     string(report.Health.Status),
		LastCheck:  report.Health.Timestamp,
		LastActive: time.Time{},
	}

	upsertContainerQuery := r.db.NewInsert().
		Model(&container).
		On("CONFLICT (docker_id) DO UPDATE").
		Set("status = EXCLUDED.status").
		Set("last_check = EXCLUDED.last_check").
		Returning("id, status, last_check as timestamp")

	_, err := r.db.NewInsert().
		With("container_data", upsertContainerQuery).
		Table("health_log", "container_data").
		Exec(ctx)

	return err
}

func (r *HealthMonitorRepositoryImpl) InsertActiveSnapshotAndUpdateContainer(report health_monitor.HealthReport, ctx context.Context) error {
	container := Container{
		DockerId:   report.DockerId,
		IpAddr:     report.IpAddr,
		Name:       report.Name,
		Status:     string(report.Health.Status),
		LastCheck:  report.Health.Timestamp,
		LastActive: report.Health.Timestamp,
	}

	upsertContainerQuery := r.db.NewInsert().
		Model(&container).
		On("CONFLICT (docker_id) DO UPDATE").
		Set("status = EXCLUDED.status").
		Set("last_check = EXCLUDED.last_check").
		Set("last_active = EXCLUDED.last_active").
		Returning("id, status, last_check as timestamp")

	_, err := r.db.NewInsert().
		With("container_data", upsertContainerQuery).
		Table("health_log", "container_data").
		Exec(ctx)

	return err
}
