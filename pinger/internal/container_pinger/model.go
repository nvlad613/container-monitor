package container_pinger

import (
	"pinger/internal/container_finder"
	"pinger/internal/container_status"
	"time"
)

type HealthReport struct {
	container_finder.Container
	Status    container_status.ContainerStatus
	Timestamp time.Time
}
