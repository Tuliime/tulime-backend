package farminputs

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateFarmInputValidator struct {
	Name          string  `validate:"string"`
	Category      string  `validate:"string"`
	Purpose       string  `validate:"string"`
	Price         float64 `validate:"number"`
	PriceCurrency string  `validate:"string"`
	Source        string  `validate:"string"`
}

var UpdateFarmInput = func(c *fiber.Ctx) error {
	farmInput := models.FarmInputs{}

	if err := c.BodyParser(&farmInput); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateFarmInputValidator
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	farmInputID := c.Params("id")

	savedFarmInput, err := farmInput.FindOne(farmInputID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFarmInput.Name = farmInput.Name
	savedFarmInput.Category = farmInput.Category
	savedFarmInput.Purpose = farmInput.Purpose
	savedFarmInput.Price = farmInput.Price
	savedFarmInput.PriceCurrency = farmInput.PriceCurrency
	savedFarmInput.Source = farmInput.Source

	updatedFarmInput, err := savedFarmInput.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedFarmInput,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
