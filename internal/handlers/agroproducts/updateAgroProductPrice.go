package agroproducts

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateAgroProductPriceInput struct {
	Amount   float64 `validate:"number"`
	Currency string  `validate:"string"`
}

var UpdateAgroProductPrice = func(c *fiber.Ctx) error {
	agroProductPrice := models.AgroproductPrice{}

	if err := c.BodyParser(&agroProductPrice); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateAgroProductPriceInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	priceID := c.Params("priceID")

	savedAgroProductPrice, err := agroProductPrice.FindOne(priceID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedAgroProductPrice.Amount = agroProductPrice.Amount
	savedAgroProductPrice.Currency = agroProductPrice.Currency

	err = savedAgroProductPrice.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    savedAgroProductPrice,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
