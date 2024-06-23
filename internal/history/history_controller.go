package history

import (
	"vitaliiPsl/synthesizer/internal/logger"
	"vitaliiPsl/synthesizer/internal/users"

	"github.com/gofiber/fiber/v2"
)

type HistoryController struct {
	service HistoryService
}

func NewHistoryController(historyService HistoryService) *HistoryController {
	return &HistoryController{service: historyService}
}

func (controller *HistoryController) HandleFetchHistory(c *fiber.Ctx) error {
	logger.Logger.Info("Handling history request...")

	tempUser := c.Locals("user")
	if tempUser == nil {
		logger.Logger.Error("No user found in context.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userDto, ok := tempUser.(*users.UserDto)
	if !ok {
		logger.Logger.Error("Failed to convert context value to UserDto")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}

	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}

	response, err := controller.service.GetHistoryRecordsByUserId(userDto, page, limit)
	if err != nil {
		logger.Logger.Error("Failed to handle history request", "message", err.Error())
		return err
	}

	logger.Logger.Info("Handled history request.")
	return c.Status(fiber.StatusOK).JSON(response)
}

func (controller *HistoryController) DeleteHistory(c *fiber.Ctx) error {
	logger.Logger.Info("Handling delete history request...")

	tempUser := c.Locals("user")
	if tempUser == nil {
		logger.Logger.Error("No user found in context.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userDto, ok := tempUser.(*users.UserDto)
	if !ok {
		logger.Logger.Error("Failed to convert context value to UserDto")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}
	controller.service.DeleteHistory(userDto.Id)

	logger.Logger.Info("Handled delete history request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (controller *HistoryController) DeleteHistoryRecord(c *fiber.Ctx) error {
	logger.Logger.Info("Handling delete history record request...")

	tempUser := c.Locals("user")
	if tempUser == nil {
		logger.Logger.Error("No user found in context.")
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Unauthorized"})
	}

	userDto, ok := tempUser.(*users.UserDto)
	if !ok {
		logger.Logger.Error("Failed to convert context value to UserDto")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	recordId := c.Params("id")
	if recordId == "" {
		logger.Logger.Error("Record Id is missing.")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Record Id is required",
		})
	}

	controller.service.DeleteHistoryRecordById(recordId, userDto.Id)

	logger.Logger.Info("Handled delete history record request.")
	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}
