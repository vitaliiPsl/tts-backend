package history

import (
	service_errors "vitaliiPsl/synthesizer/internal/error"
	"vitaliiPsl/synthesizer/internal/logger"

	"github.com/gofiber/fiber/v2"
)

type HistoryController struct {
	service *HistoryService
}

func NewHistoryController(historyService *HistoryService) *HistoryController {
	return &HistoryController{service: historyService}
}

func (controller *HistoryController) HandleFetchHistory(c *fiber.Ctx) error {
	logger.Logger.Info("Handling history request...")

	var userId string
	var ok bool

	userIdInterface := c.Locals("userId")
	if userIdInterface == nil {
		logger.Logger.Error("User Id is not present in context")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{})
	}

	userId, ok = userIdInterface.(string)
	if !ok {
		logger.Logger.Error("User Id is not a string.")
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

	response, err := controller.service.GetHistoryRecordsByUserId(userId, page, limit)
	if err != nil {
		logger.Logger.Error("Failed to handle history request", "message", err.Error())
		return service_errors.HandleError(c, err)
	}

	logger.Logger.Info("Handled history request.")
	return c.Status(fiber.StatusOK).JSON(response)
}
