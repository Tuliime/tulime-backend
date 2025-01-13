package farmmanager

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var DeleteFarmManager = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	farmManagerID := c.Params("id")

	savedFarmManager, err := farmManager.FindOne(farmManagerID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmManager.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Farm manager of provided id is not found!")
	}

	if err := farmManager.Delete(farmManagerID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Farm Manager deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
