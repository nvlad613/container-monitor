package health_monitor

import "context"

func (s *ServiceImpl) GetContainersList(offset int, limit int, ctx context.Context) ([]Container, error) {
	// TODO handle errors
	return s.repository.GetContainers(offset, limit, ctx)
}
