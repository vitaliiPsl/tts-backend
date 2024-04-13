package router

import (
	"vitaliiPsl/synthesizer/internal/auth"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App, authController *auth.AuthController) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	app.Use(cors.New())

	api := app.Group("/v1")

	// Auth
	auth := api.Group("/auth")
	auth.Post("/sign-up", authController.HandleSignUp)
	auth.Post("/sign-in", authController.HandleSignIn)
	auth.Post("/verify-email", authController.HandleEmailVerification)

}
