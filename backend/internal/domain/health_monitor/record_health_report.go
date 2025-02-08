package health_monitor

import (
	"backend/internal/domain"
	"backend/internal/domain/container_status"
	"context"
	"fmt"
)

func (s *ServiceImpl) RecordHealthReport(report HealthReport, ctx context.Context) error {
	var err error
	switch report.Health.Status {
	case container_status.Offline:
		err = s.repository.InsertSnapshotAndUpdateContainer(report, ctx)
	case container_status.Online:
		err = s.repository.InsertActiveSnapshotAndUpdateContainer(report, ctx)
	default:
		return fmt.Errorf("%w: %s", domain.UnknownContainerStatusError, report.Health.Status)
	}

	// TODO: handle errors
	return err
}
