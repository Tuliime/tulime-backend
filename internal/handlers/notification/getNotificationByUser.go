package notification

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetNotificationByUser = func(c *fiber.Ctx) error {
	notification := models.Notification{}
	userID := c.Params("userID")
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")

	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID")
	}

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}

	notifications, err := notification.FindByUser(userID, limit, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	allNotificationCount, err := notification.FindUnreadCountByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	chatNotificationCount, err := notification.FindUnreadCountByUserAndType(userID, "chat")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pagination := fiber.Map{
		"limit":  limit,
		"cursor": cursorParam,
	}

	data := fiber.Map{
		"allNotificationCount":  allNotificationCount,
		"chatNotificationCount": chatNotificationCount,
		"notifications":         notifications,
		"pagination":            pagination,
	}

	response := fiber.Map{
		"status": "success",
		"data":   data,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
