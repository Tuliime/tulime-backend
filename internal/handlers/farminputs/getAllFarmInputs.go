package farminputs

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllFarmInputs = func(c *fiber.Ctx) error {
	farmInputs := models.FarmInputs{}
	limitParam := c.Query("limit")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	allFarmInputs, err := farmInputs.FindAll(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pagination := map[string]interface{}{
		"limit": limit,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allFarmInputs,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
