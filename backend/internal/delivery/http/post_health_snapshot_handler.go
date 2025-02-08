package http

import (
	"backend/internal/delivery"
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

// @Summary		Add health snapshot to the container health history
// @Description	Add health snapshot to the container health history and upsert container data
// @Tags			containers
// @Accept			json
// @Produce			json
// @Param			request		body	delivery.ContainerHealthReportView	true	"Snapshot and container data"
// @Success		200
// @Failure		400	{object}	Error
// @Failure		500	{object}	Error
// @Router			/container/health [post]
func (s *Server) postHealthSnapshotHandler(ctx *fiber.Ctx) error {
	var requestBody delivery.ContainerHealthReportView
	if err := ctx.BodyParser(&requestBody); err != nil {
		s.logger.Warn(err)

		return SendError(fiber.StatusBadRequest, "invalid request body", ctx)
	}

	if err := s.validation.Struct(&requestBody); err != nil {
		var invalidValidationError *validator.InvalidValidationError
		if errors.As(err, &invalidValidationError) {
			s.logger.Error(err)

			return SendError(fiber.StatusInternalServerError, "server error", ctx)
		}

		return SendError(fiber.StatusBadRequest, err.Error(), ctx)
	}

	reportModel, err := requestBody.ToModel()
	if err != nil {
		return SendError(fiber.StatusInternalServerError, "invalid timestamp format", ctx)
	}

	if err = s.healthMonitorService.RecordHealthReport(reportModel, ctx.Context()); err != nil {
		s.logger.Error(err)

		return SendError(fiber.StatusInternalServerError, "server error", ctx)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}
