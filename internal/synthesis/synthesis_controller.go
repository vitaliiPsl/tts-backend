package synthesis

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/user"
	"vitaliiPsl/synthesizer/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type SynthesisController struct {
	synthesisService  *SynthesisService
	validationService *validation.ValidationService
}

func NewSynthesisController(synthesisService *SynthesisService, validationService *validation.ValidationService) *SynthesisController {
	return &SynthesisController{
		synthesisService:  synthesisService,
		validationService: validationService,
	}
}

func (controller *SynthesisController) HandleSynthesis(c *fiber.Ctx) error {
	logger.Logger.Info("Handling speech synthesis...")

	var ok bool
	var userDto *user.UserDto

	tempUser := c.Locals("user")
	if tempUser != nil {
		userDto, ok = tempUser.(*user.UserDto)
		if !ok {
			logger.Logger.Error("Failed to convert context value to UserDto")
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
		}
	}

	userId := ""
	if userDto != nil {
		userId = userDto.Id
	}

	var req requests.SynthesisRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse synthesis request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateSynthesisRequest(&req); err != nil {
		logger.Logger.Error("Synthesis request didn't pass validation", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	result, err := controller.synthesisService.HandleSynthesisRequest(&req, userId)
	if err != nil {
		logger.Logger.Error("Failed to synthesize speech", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled speech synthesis.")
	return c.Status(fiber.StatusOK).JSON(result)
}
