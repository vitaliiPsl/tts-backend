package auth

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	authService       *AuthService
	validationService *validation.ValidationService
}

func NewAuthController(authService *AuthService, validationService *validation.ValidationService) *AuthController {
	return &AuthController{authService: authService, validationService: validationService}
}

func (controller *AuthController) HandleSignUp(c *fiber.Ctx) error {
	logger.Logger.Info("Handling sign up request...")

	var req requests.SignUpRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse sign up request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateSignUpRequest(&req); err != nil {
		logger.Logger.Error("Sign up request didn't pass validation", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	if err := controller.authService.HandleSignUp(&req); err != nil {
		logger.Logger.Error("Failed to handle sign up request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled sign up request.")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Created"})
}

func (controller *AuthController) HandleSignIn(c *fiber.Ctx) error {
	logger.Logger.Info("Handling sign in request...")

	var req requests.SignInRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse sign up request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateSignInRequest(&req); err != nil {
		logger.Logger.Error("Sign in request didn't pass validation", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	token, err := controller.authService.HandleSignIn(&req)
	if err != nil {
		logger.Logger.Error("Failed to handle sign in request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled sign in request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": token})
}

func (controller *AuthController) HandleEmailVerification(c *fiber.Ctx) error {
	logger.Logger.Info("Handling email verification request...")

	var req requests.EmailVerificationRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse email verification request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateEmailVerificationRequest(&req); err != nil {
		logger.Logger.Error("Password reset request didn't pass validation", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	err := controller.authService.HandleEmailVerification(&req)
	if err != nil {
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled email verification.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
