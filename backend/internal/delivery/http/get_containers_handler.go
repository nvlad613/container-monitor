package http

import (
	"backend/internal/delivery"
	"backend/internal/domain/health_monitor"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"strconv"
)

// @Summary		Get list of containers
// @Description	Get list of containers with offset and limit parameters
// @Tags			containers
// @Produce		json
// @Param			offset	query		int	false	"Offset"
// @Param			limit	query		int	false	"Limit"
// @Success		200	{object}	[]delivery.ContainerView
// @Failure		400	{object}	Error
// @Failure		500	{object}	Error
// @Router			/container [get]
func (s *Server) getContainersHandler(ctx *fiber.Ctx) error {
	offsetRaw := ctx.Query("offset", "0")
	limitRaw := ctx.Query("limit", "0")

	var (
		err    error
		offset int
		limit  int
	)
	if offset, err = strconv.Atoi(offsetRaw); err != nil || offset < 0 {
		return SendError(fiber.StatusBadRequest, "invalid offset value", ctx)
	}
	if limit, err = strconv.Atoi(limitRaw); err != nil || limit < 0 {
		return SendError(fiber.StatusBadRequest, "invalid limit value", ctx)
	}

	containers, err := s.healthMonitorService.GetContainersList(offset, limit, ctx.Context())
	if err != nil {
		s.logger.Error(err)

		return SendError(fiber.StatusInternalServerError, "server error", ctx)
	}

	containerViews := lo.Map(containers, func(item health_monitor.Container, _ int) delivery.ContainerView {
		return new(delivery.ContainerView).FromModel(item)
	})

	return ctx.JSON(containerViews)
}
