package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateAdvertPrice = func(c *fiber.Ctx) error {
	advertPriceID := c.Params("advertPriceID")
	advertPrice := models.AdvertPrice{}

	savedAdvertPrice, err := advertPrice.FindOne(advertPriceID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvertPrice.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert price of provided id doesn't exist!")
	}

	if err := c.BodyParser(&advertPrice); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advertPrice.Amount == 0 || advertPrice.Currency == "" || advertPrice.Unit == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Amount/Currency/Unit!")
	}

	savedAdvertPrice.Amount = advertPrice.Amount
	savedAdvertPrice.Currency = advertPrice.Currency
	savedAdvertPrice.Unit = advertPrice.Unit

	updatedAdvertPrice, err := savedAdvertPrice.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert price updated successfully!",
		"data":    updatedAdvertPrice,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
