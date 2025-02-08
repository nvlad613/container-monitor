package health_monitor

import (
	"context"
	"time"
)

func (s *ServiceImpl) GetContainerHealth(containerId int, from time.Time, to time.Time, ctx context.Context) (*HealthLog, error) {
	//TODO handle errors
	snapshots, err := s.repository.GetHealthSnapshots(containerId, from, to, ctx)
	if err != nil {
		return nil, err
	}

	return &HealthLog{
		ContainerId: containerId,
		FromDate:    from,
		ToDate:      to,
		Log:         snapshots,
	}, nil
}
