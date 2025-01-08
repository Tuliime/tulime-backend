package agroproducts

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type PostAgroProductPriceInput struct {
	AgroproductID string `validate:"string"`
	Amount        string `validate:"number"`
	Currency      string `validate:"string"`
}

var PostAgroProductPrice = func(c *fiber.Ctx) error {
	agroProductPrice := models.AgroproductPrice{}

	if err := c.BodyParser(agroProductPrice); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input PostAgroProductPriceInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	newAgroProductPrice, err := agroProductPrice.Create(agroProductPrice)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Created successfully!",
		"data":    newAgroProductPrice,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
