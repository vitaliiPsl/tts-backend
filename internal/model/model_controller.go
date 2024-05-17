package model

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/requests"
	"vitaliiPsl/synthesizer/internal/validation"

	"github.com/gofiber/fiber/v2"
)

type ModelController struct {
	service           *ModelService
	validationService *validation.ValidationService
}

func NewModelController(modelService *ModelService, validationService *validation.ValidationService) *ModelController {
	return &ModelController{service: modelService, validationService: validationService}
}

func (controller *ModelController) HandleSaveModel(c *fiber.Ctx) error {
	logger.Logger.Info("Handling save model request...")

	var req requests.ModelRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse model request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if err := controller.validationService.ValidateModelRequest(&req); err != nil {
		logger.Logger.Error("Model request didn't pass validation", "message", err.Error(), err)
		return service_errors.HandleError(c, err)
	}

	if _, err := controller.service.SaveModel(&req); err != nil {
		logger.Logger.Error("Failed to handle save model request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled save model request.")
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{})
}

func (controller *ModelController) HandleUpdateModel(c *fiber.Ctx) error {
	logger.Logger.Info("Handling update model request...")

	modelId := c.Params("id")
	if modelId == "" {
		logger.Logger.Error("Model Id is missing.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Model Id is required",
		})
	}

	var req requests.ModelRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Logger.Error("Failed to parse model request", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	if _, err := controller.service.UpdateModel(modelId, &req); err != nil {
		logger.Logger.Error("Failed to handle update model request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled update model request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (controller *ModelController) HandleDeleteModel(c *fiber.Ctx) error {
	logger.Logger.Info("Handling delete model request...")

	modelId := c.Params("id")
	if modelId == "" {
		logger.Logger.Error("Model Id is missing.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Model Id is required",
		})
	}

	if err := controller.service.DeleteModel(modelId); err != nil {
		logger.Logger.Error("Failed to delete model", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled delete model request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (controller *ModelController) HandleFetchModels(c *fiber.Ctx) error {
	logger.Logger.Info("Handling models request...")

	response, err := controller.service.GetModels()
	if err != nil {
		logger.Logger.Error("Failed to handle models request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled models request.")
	return c.Status(fiber.StatusOK).JSON(response)
}
