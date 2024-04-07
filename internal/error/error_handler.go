package error

import "github.com/gofiber/fiber/v2"

func HandleError(ctx *fiber.Ctx, err error) error {
	switch e := err.(type) {
	case *ErrNotFound:
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *ErrBadRequest:
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *ErrUnauthorized:
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *ErrInternalServer:
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": e.Error(),
		})
	case *ErrBadGateway:
		return ctx.Status(fiber.StatusBadGateway).JSON(fiber.Map{
			"error": e.Error(),
		})
	default:
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
}
