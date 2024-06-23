package server

import (
	internal_errors "vitaliiPsl/synthesizer/internal/error"

	"github.com/gofiber/fiber/v2"
)

type FiberServer struct {
	*fiber.App
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ErrorHandler: internal_errors.ErrorHandler,
		}),
	}

	return server
}
