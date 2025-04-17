package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PostAdvertPrice = func(c *fiber.Ctx) error {
	advertID := c.Params("id")
	advert := models.Advert{}
	advertPrice := models.AdvertPrice{}

	savedAdvert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	savedAdvertPrice, err := advertPrice.FindByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvertPrice.ID != "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id already has a price!")
	}

	advertPrice.AdvertID = advertID

	if err := c.BodyParser(&advertPrice); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advertPrice.Amount == 0 || advertPrice.Currency == "" || advertPrice.Unit == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Amount/Currency/Unit!")
	}

	newAdvertPrice, err := advertPrice.Create(advertPrice)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert price created successfully!",
		"data":    newAdvertPrice,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
