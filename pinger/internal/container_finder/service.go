package container_finder

import (
	"context"
	"errors"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/go-playground/validator/v10"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type Finder interface {
	ContainersList(ctx context.Context) ([]Container, error)
}

type FinderImpl struct {
	logger       *zap.SugaredLogger
	dockerClient *client.Client
	validate     *validator.Validate
}

func New(
	logger *zap.SugaredLogger,
	dockerClient *client.Client,
	validate *validator.Validate,
) *FinderImpl {
	return &FinderImpl{
		logger:       logger,
		dockerClient: dockerClient,
		validate:     validate,
	}
}

func (f *FinderImpl) ContainersList(ctx context.Context) ([]Container, error) {
	containers, err := f.dockerClient.ContainerList(ctx, container.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get list of containers: %w", err)
	}

	containerAddresses := make([]Container, 0, len(containers))
	for _, cnt := range containers {
		inspect, err := f.dockerClient.ContainerInspect(ctx, cnt.ID)
		if err != nil {
			f.logger.Errorw("failed to inspect container: "+err.Error(),
				"id", cnt.ID, "name", cnt.Names[0])
			continue
		}

		cntName, exist := lo.First(cnt.Names)
		if !exist {
			cntName = "unnamed"
		}

		var (
			ip   = inspect.NetworkSettings.IPAddress
			port string
		)
		if err = f.validate.Var(&ip, "required,ip4_addr"); err != nil {
			f.logger.Warnw("container have no ip address in default network; trying to find mapped port..",
				"id", cnt.ID, "name", cntName)

			// Checking mapped ports
			if ip, port, err = extractPortBinding(&inspect); err != nil {
				f.logger.Errorw("failed to get container address: "+err.Error(),
					"id", cnt.ID, "name", cntName)
				continue
			}
		} else {
			if port, err = extractExposedPort(&inspect); err != nil {
				f.logger.Errorw("failed to get container address: "+err.Error(),
					"id", cnt.ID, "name", cntName)
				continue
			}
		}

		containerAddresses = append(containerAddresses, Container{
			Name:     cntName,
			DockerId: cnt.ID,
			IpAddr:   ip,
			Port:     port,
		})
	}

	return containerAddresses, nil
}

func extractExposedPort(inspect *types.ContainerJSON) (string, error) {
	for portObj, _ := range inspect.Config.ExposedPorts {
		if portObj.Proto() != "tcp" {
			continue
		}

		return portObj.Port(), nil
	}

	return "", errors.New("exposed port not found")
}

func extractPortBinding(inspect *types.ContainerJSON) (string, string, error) {
	bindings := inspect.HostConfig.PortBindings
	for _, hostAddrs := range bindings {
		addr, exist := lo.First(hostAddrs)
		if !exist {
			continue
		}

		ip := addr.HostIP
		port := addr.HostPort
		if ip == "" {
			ip = "127.0.0.1"
		}

		return ip, port, nil
	}

	return "", "", errors.New("port bindings not found")
}
