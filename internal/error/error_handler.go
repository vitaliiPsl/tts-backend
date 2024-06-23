package internal_errors

import (
	"github.com/gofiber/fiber/v2"
	"vitaliiPsl/synthesizer/internal/logger"
)

func ErrorHandler(ctx *fiber.Ctx, err error) error {
	logger.Logger.Error("Error while handling request", "error", err.Error())

	switch e := err.(type) {
	case *ErrNotFound:
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": e.Error()})
	case *ErrBadRequest:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": e.Error()})
	case *ErrForbidden:
		return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": e.Error()})
	case *ErrUnauthorized:
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": e.Error()})
	default:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
}
