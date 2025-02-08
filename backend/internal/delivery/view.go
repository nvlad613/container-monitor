package delivery

import (
	"backend/internal/domain/container_status"
	"backend/internal/domain/health_monitor"
	"github.com/samber/lo"
	"time"
)

type ContainerView struct {
	Name         string                           `json:"name"`
	Id           int                              `json:"id"`
	IpAddr       string                           `json:"ip"`
	Status       container_status.ContainerStatus `json:"status"`
	LastCheck    time.Time                        `json:"last_check"`
	LastActivity time.Time                        `json:"last_activity"`
}

func (view *ContainerView) FromModel(model health_monitor.Container) ContainerView {
	*view = ContainerView{
		Name:         model.Name,
		Id:           model.Id,
		IpAddr:       model.IpAddr,
		Status:       model.Status,
		LastCheck:    model.LastCheck,
		LastActivity: model.LastActivity,
	}

	return *view
}

type ContainerHealthReportView struct {
	Name     string         `json:"name" validate:"required,max=64"`
	DockerId string         `json:"docker_id" validate:"required,len=64"`
	IpAddr   string         `json:"ip" validate:"required,ipv4"`
	Health   HealthSnapshot `json:"health_report" validate:"required"`
}

func (view ContainerHealthReportView) ToModel() (health_monitor.HealthReport, error) {
	parsedTime, err := time.Parse(time.DateTime, view.Health.Timestamp)
	if err != nil {
		return health_monitor.HealthReport{}, err
	}

	return health_monitor.HealthReport{
		Name:     view.Name,
		DockerId: view.DockerId,
		IpAddr:   view.IpAddr,
		Health: health_monitor.HealthSnapshot{
			Status:    container_status.ContainerStatus(view.Health.Status),
			Timestamp: parsedTime,
		},
	}, nil
}

type HealthLogView struct {
	ContainerId int              `json:"container_id"`
	FromDate    time.Time        `json:"from_date"`
	ToDate      time.Time        `json:"to_date"`
	HealthLog   []HealthSnapshot `json:"health_log"`
}

func (view *HealthLogView) FromModel(model health_monitor.HealthLog) HealthLogView {
	*view = HealthLogView{
		ContainerId: model.ContainerId,
		FromDate:    model.FromDate,
		ToDate:      model.ToDate,
		HealthLog: lo.Map(model.Log, func(item health_monitor.HealthSnapshot, _ int) HealthSnapshot {
			return HealthSnapshot{
				Status:    string(item.Status),
				Timestamp: item.Timestamp.Format(time.DateTime),
			}
		}),
	}

	return *view
}

type HealthSnapshot struct {
	Status    string `json:"status" validate:"required,oneof=offline online"`
	Timestamp string `json:"timestamp" validate:"required"`
}
