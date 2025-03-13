package notification

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateNotificationAsRead = func(c *fiber.Ctx) error {
	notification := models.Notification{}
	notificationID := c.Params("id")

	if notificationID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide notificationID!")
	}

	savedNotification, err := notification.FindOne(notificationID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedNotification.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Notification of provided id does'nt exist!")
	}

	if savedNotification.IsRead {
		return fiber.NewError(fiber.StatusBadRequest, "Notification is already updated as read!")
	}

	savedNotification.IsRead = true

	updateNotification, err := savedNotification.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Notification updated successfully!",
		"data":    updateNotification,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
