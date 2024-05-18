package router

import (
	"vitaliiPsl/synthesizer/internal/auth"
	"vitaliiPsl/synthesizer/internal/history"
	"vitaliiPsl/synthesizer/internal/model"
	"vitaliiPsl/synthesizer/internal/synthesis"
	"vitaliiPsl/synthesizer/internal/user"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(
	app *fiber.App,
	authMiddleware *auth.AuthMiddleware,
	authController *auth.AuthController,
	modelController *model.ModelController,
	synthesisController *synthesis.SynthesisController,
	historyController *history.HistoryController,
) {

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK"})
	})

	app.Use(cors.New())

	api := app.Group("/v1")

	// Auth
	authApi := api.Group("/auth")
	authApi.Get("/me", authMiddleware.ProtectedRoute(), authController.HandleAuthenticatedUserRequest)
	authApi.Post("/sign-up", authController.HandleSignUp)
	authApi.Post("/sign-in", authController.HandleSignIn)
	authApi.Get("/sso/:provider", authController.HandleSsoSignIn)
	authApi.Post("/sso/:provider/sign-in", authController.HandleSsoCallback)
	authApi.Post("/verify-email", authController.HandleEmailVerification)
	authApi.Post("/reset-password", authController.HandleResetPassword)
	authApi.Post("/send-password-reset-email", authController.HandleSendPasswordResetToken)

	modelApi := api.Group("/models")
	modelApi.Post("", authMiddleware.ProtectedRoute(user.RoleAdmin), modelController.HandleSaveModel)
	modelApi.Patch(":id", authMiddleware.ProtectedRoute(user.RoleAdmin), modelController.HandleUpdateModel)
	modelApi.Delete(":id", authMiddleware.ProtectedRoute(user.RoleAdmin), modelController.HandleDeleteModel)
	modelApi.Get("", authMiddleware.OpenRoute(), modelController.HandleFetchModels)

	synthesisApi := api.Group("/synthesis")
	synthesisApi.Post("", authMiddleware.OpenRoute(), synthesisController.HandleSynthesis)

	historyApi := api.Group("/history")
	historyApi.Get("", authMiddleware.ProtectedRoute(), historyController.HandleFetchHistory)
	historyApi.Delete("", authMiddleware.ProtectedRoute(), historyController.DeleteHistory)
	historyApi.Delete(":id", authMiddleware.ProtectedRoute(), historyController.DeleteHistoryRecord)
}
