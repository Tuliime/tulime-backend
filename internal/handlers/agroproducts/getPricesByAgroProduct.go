package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetPricesByAgroProduct = func(c *fiber.Ctx) error {
	agroProductPrice := models.AgroproductPrice{}
	agroProductID := c.Params("id")

	agroProductPrices, err := agroProductPrice.FindByAgroProduct(agroProductID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   agroProductPrices,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
