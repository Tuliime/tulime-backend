package chatroom

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetOnlineStatus = func(c *fiber.Ctx) error {
	onlineStatus := models.OnlineStatus{}

	onlineStatuses, err := onlineStatus.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   onlineStatuses,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
