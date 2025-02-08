package http

import (
	"backend/internal/delivery"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

const (
	fromDateDefault = "1980-01-01 00:00:00"
)

// @Summary		Get history of container statuses
// @Description	Get history of certain container health between from_date and to_date
// @Tags		containers
// @Produce			json
// @Param			from_date	query		string	false	"Start of period"
// @Param			to_date		query		string	false	"End of period"
// @Param			id			path		int		true	"Container id"
// @Success		200	{object}	[]delivery.HealthLogView
// @Failure		400	{object}	Error
// @Failure		500	{object}	Error
// @Router			/container/{id} [get]
func (s *Server) getHealthLogHandler(ctx *fiber.Ctx) error {
	var (
		containerIdRaw = ctx.Params("container_id")
		fromDateRaw    = ctx.Query("from_date", fromDateDefault)
		toDateRaw      = ctx.Query("to_date", time.Now().Format(time.DateTime))
	)

	var (
		err         error
		containerId int
		fromDate    time.Time
		toDate      time.Time
	)
	if containerId, err = strconv.Atoi(containerIdRaw); err != nil || containerId < 1 {
		return SendError(fiber.StatusBadRequest, "invalid container id", ctx)
	}
	if fromDate, err = time.Parse(time.DateTime, fromDateRaw); err != nil {
		return SendError(fiber.StatusBadRequest, "invalid start of the period", ctx)
	}
	if toDate, err = time.Parse(time.DateTime, toDateRaw); err != nil {
		return SendError(fiber.StatusBadRequest, "invalid end of the period", ctx)
	}

	containerHealth, err := s.healthMonitorService.GetContainerHealth(
		containerId,
		fromDate,
		toDate,
		ctx.Context(),
	)
	if err != nil {
		s.logger.Error(err)

		return SendError(fiber.StatusInternalServerError, "server error", ctx)
	}

	return ctx.JSON(new(delivery.HealthLogView).FromModel(*containerHealth))
}
