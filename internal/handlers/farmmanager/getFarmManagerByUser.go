package farmmanager

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetFarmManagerByUser = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	userID := c.Params("userID")

	farmManager, err := farmManager.FindByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   farmManager,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
