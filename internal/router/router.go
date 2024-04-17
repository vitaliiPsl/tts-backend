package router

import (
	"vitaliiPsl/synthesizer/internal/auth"
	"vitaliiPsl/synthesizer/internal/auth/jwt"
	"vitaliiPsl/synthesizer/internal/synthesis"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(
	app *fiber.App,
	authController *auth.AuthController,
	synthesisController *synthesis.SynthesisController,
	jwtService *jwt.JwtService,
) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	app.Use(cors.New())

	api := app.Group("/v1")

	// Auth
	authApi := api.Group("/auth")
	authApi.Post("/sign-up", authController.HandleSignUp)
	authApi.Post("/sign-in", authController.HandleSignIn)
	authApi.Get("/sso/:provider", authController.HandleSsoSignIn)
	authApi.Post("/sso/:provider/sign-in", authController.HandleSsoCallback)
	authApi.Post("/verify-email", authController.HandleEmailVerification)
	authApi.Post("/reset-password", authController.HandleResetPassword)
	authApi.Post("/send-password-reset-email", authController.HandleSendPasswordResetToken)

	synthesisApi := api.Group("/synthesis")
	synthesisApi.Post("", synthesisController.HandleSynthesis)

}
