package notification

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetNotificationByUser = func(c *fiber.Ctx) error {
	notification := models.Notification{}
	userID := c.Params("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID")
	}

	notifications, err := notification.FindUnreadByUser(userID)
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

	data := fiber.Map{
		"allNotificationCount":  allNotificationCount,
		"chatNotificationCount": chatNotificationCount,
		"notifications":         notifications,
	}

	response := fiber.Map{
		"status": "success",
		"data":   data,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
