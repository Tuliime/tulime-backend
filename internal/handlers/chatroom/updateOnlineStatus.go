package chatroom

import (
	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateOnlineStatus = func(c *fiber.Ctx) error {
	onlineStatus := models.OnlineStatus{}

	if err := c.BodyParser(&onlineStatus); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if onlineStatus.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID!")
	}

	updatedOnlineStatus, err := onlineStatus.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	events.EB.Publish("onlineStatus", updatedOnlineStatus)

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedOnlineStatus,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
