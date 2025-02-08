package repository

import (
	"backend/internal/domain/container_status"
	"backend/internal/domain/health_monitor"
	"github.com/uptrace/bun"
	"time"
)

type Container struct {
	bun.BaseModel `bun:"table:containers"`

	Id         int       `bun:"id,pk,autoincrement"`
	DockerId   string    `bun:"docker_id,type:char(64),notnull"`
	IpAddr     string    `bun:"ip_addr,type:inet"`
	Name       string    `bun:"name,type:varchar(128),notnull"`
	Status     string    `bun:"status,type:container_status,default:'offline'"`
	LastCheck  time.Time `bun:"last_check,type:timestamp"`
	LastActive time.Time `bun:"last_active,type:timestamp"`
}

func (entity Container) ToModel() health_monitor.Container {
	return health_monitor.Container{
		Name:         entity.Name,
		Id:           entity.Id,
		IpAddr:       entity.IpAddr,
		Status:       container_status.ContainerStatus(entity.Status),
		LastCheck:    entity.LastCheck,
		LastActivity: entity.LastActive,
	}
}

type HealthSnapshot struct {
	bun.BaseModel `bun:"table:health_log"`

	ContainerId int       `bun:"container_id,notnull"`
	Status      string    `bun:"status,type:container_status,default:'offline'"`
	Timestamp   time.Time `bun:"timestamp,notnull"`
}

func (entity HealthSnapshot) ToModel() health_monitor.HealthSnapshot {
	return health_monitor.HealthSnapshot{
		Status:    container_status.ContainerStatus(entity.Status),
		Timestamp: entity.Timestamp,
	}
}
