package api_client

import (
	"pinger/internal/container_pinger"
	"time"
)

type HealthReport struct {
	Name     string         `json:"name"`
	DockerId string         `json:"docker_id"`
	IpAddr   string         `json:"ip"`
	Port     string         `json:"port"`
	Health   HealthSnapshot `json:"health_report"`
}

type HealthSnapshot struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func (view *HealthReport) FromModel(model container_pinger.HealthReport) HealthReport {
	*view = HealthReport{
		Name:     model.Name,
		DockerId: model.DockerId,
		IpAddr:   model.IpAddr,
		Port:     model.Port,
		Health: HealthSnapshot{
			Status:    string(model.Status),
			Timestamp: model.Timestamp.Format(time.DateTime),
		},
	}

	return *view
}
