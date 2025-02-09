package container_pinger

import (
	"fmt"
	"net"
	"pinger/config"
	"pinger/internal/container_finder"
	"pinger/internal/container_status"
	"time"
)

type Pinger interface {
	Ping(container container_finder.Container) HealthReport
}

type PingerImpl struct {
	conf config.PingerParams
}

func New(
	pingerConfig config.PingerParams,
) *PingerImpl {
	return &PingerImpl{
		conf: pingerConfig,
	}
}

func (p PingerImpl) Ping(container container_finder.Container) HealthReport {
	reqTimeout := time.Duration(p.conf.PingPeriodSec) * time.Second
	timestamp := time.Now()
	conn, err := net.DialTimeout(
		"tcp",
		fmt.Sprintf("%s:%s", container.IpAddr, container.Port),
		reqTimeout,
	)
	if err != nil {
		return HealthReport{
			Container: container,
			Status:    container_status.Offline,
			Timestamp: timestamp,
		}
	}
	conn.Close()

	return HealthReport{
		Container: container,
		Status:    container_status.Online,
		Timestamp: timestamp,
	}
}
