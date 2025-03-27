package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PostAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	if err := c.BodyParser(&advert); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advert.StoreID == "" || advert.UserID == "" ||
		advert.ProductName == "" || advert.ProductDescription == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"Missing StoreID/UserID/ProductName/ProductDescription!")
	}

	newStore, err := advert.Create(advert)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert created successfully!",
		"data":    newStore,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
