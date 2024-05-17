package auth

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/user"
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

func (controller *AuthController) HandleAuthenticatedUserRequest(c *fiber.Ctx) error {
	logger.Logger.Info("Handling authenticated user request...")

	tempUser := c.Locals("user")
	if tempUser == nil {
		logger.Logger.Error("No user found in context.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userDto, ok := tempUser.(*user.UserDto)
	if !ok {
		logger.Logger.Error("Failed to convert context value to UserDto")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	logger.Logger.Info("Handled authenticated user request.")
	return c.Status(fiber.StatusOK).JSON(userDto)
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


func (controller *AuthController) HandleSsoSignIn(c *fiber.Ctx) error {
	logger.Logger.Info("Handling SSO sign in request...")

	provider := c.Params("provider")
	if provider == "" {
		logger.Logger.Warn("SSO provider is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SSO provider is required."})
	}

	url, err := controller.authService.HandleSsoSignIn(provider)
	if err != nil {
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled SSO sign in request. Redirecting to SSO provider", "url", url)
	return c.Redirect(url, fiber.StatusFound)
}

func (controller *AuthController) HandleSsoCallback(c *fiber.Ctx) error {
	logger.Logger.Info("Handling SSO callback request...")

	provider := c.Params("provider")
	if provider == "" {
		logger.Logger.Warn("SSO provider is missing")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "SSO provider is required."})
	}

	var req requests.SignInWithSSORequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse SSO callback request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	jwtToken, err := controller.authService.HandleSSOCallback(provider, req.Code)
	if err != nil {
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled SSO callback request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"token": jwtToken})
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

func (controller *AuthController) HandleSendPasswordResetToken(c *fiber.Ctx) error {
	logger.Logger.Info("Handling 'send password reset token' request...")

	var req requests.VerificationTokenRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse 'send password verification token' request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateVerificationTokenRequest(&req); err != nil {
		logger.Logger.Error("'Send password reset' request didn't pass validation", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	err := controller.authService.HandleSendPasswordResetToken(&req)
	if err != nil {
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled 'send password reset token' request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (controller *AuthController) HandleResetPassword(c *fiber.Ctx) error {
	logger.Logger.Info("Handling password reset request...")

	var req requests.PasswordResetRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse password reset request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidatePasswordResetRequest(&req); err != nil {
		logger.Logger.Error("Password reset request didn't pass validation", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	err := controller.authService.HandlePasswordReset(&req)
	if err != nil {
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled password reset request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}