package farmmanager

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetFarmManager = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	farmManagerID := c.Params("id")

	farmManager, err := farmManager.FindOne(farmManagerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   farmManager,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
