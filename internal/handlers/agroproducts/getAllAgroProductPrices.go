package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetAllAgroProductPrices = func(c *fiber.Ctx) error {
	agroProductPrice := models.AgroproductPrice{}

	agroProductPrices, err := agroProductPrice.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   agroProductPrices,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
