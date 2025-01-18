package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var DeleteAgroProductPrice = func(c *fiber.Ctx) error {
	agroProductPrice := models.AgroproductPrice{}

	priceID := c.Params("priceID")

	savedAgroProductPrice, err := agroProductPrice.FindOne(priceID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedAgroProductPrice.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Price of provided id is not found!")
	}

	err = agroProductPrice.Delete(priceID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Price deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
