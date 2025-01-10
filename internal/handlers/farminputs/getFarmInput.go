package farminputs

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetFarmInput = func(c *fiber.Ctx) error {
	farmInputs := models.FarmInputs{}
	farmInputID := c.Params("id")

	farmInput, err := farmInputs.FindOne(farmInputID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   farmInput,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
