package health_monitor

import (
	"backend/internal/domain/container_status"
	"time"
)

type Container struct {
	Name         string
	Id           int
	IpAddr       string
	Status       container_status.ContainerStatus
	LastCheck    time.Time
	LastActivity time.Time
}

type HealthReport struct {
	Name     string
	DockerId string
	IpAddr   string
	Health   HealthSnapshot
}

type HealthLog struct {
	ContainerId int
	FromDate    time.Time
	ToDate      time.Time
	Log         []HealthSnapshot
}

type HealthSnapshot struct {
	Status    container_status.ContainerStatus
	Timestamp time.Time
}
