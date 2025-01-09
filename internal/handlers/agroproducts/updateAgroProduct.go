package agroproducts

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateAgroProductInput struct {
	Name     string `validate:"string"`
	Category string `validate:"string"`
}

var UpdateAgroProduct = func(c *fiber.Ctx) error {
	agroProduct := models.Agroproduct{}

	if err := c.BodyParser(&agroProduct); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateAgroProductInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	agroProductID := c.Params("id")

	savedAgroProduct, err := agroProduct.FindOne(agroProductID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedAgroProduct.Name = agroProduct.Name
	savedAgroProduct.Category = agroProduct.Category

	updatedAgroProduct, err := savedAgroProduct.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedAgroProduct,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
