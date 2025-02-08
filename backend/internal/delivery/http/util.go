package http

import (
	"github.com/gofiber/fiber/v2"
)

type Error struct {
	Status     string `json:"status"`
	StatusCode int    `json:"status_code"`
	Details    string `json:"details"`
}

func SendError(status int, message string, ctx *fiber.Ctx) error {
	ctx = ctx.Status(status)
	return ctx.JSON(&Error{
		Status:     "error",
		StatusCode: status,
		Details:    message,
	}, "application/problem+json")
}
