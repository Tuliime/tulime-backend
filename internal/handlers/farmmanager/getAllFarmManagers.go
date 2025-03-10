package farmmanager

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllFarmManagers = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	limitParam := c.Query("limit")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	farmManagers, err := farmManager.FindAll(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pagination := map[string]interface{}{
		"limit": limit,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       farmManagers,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
