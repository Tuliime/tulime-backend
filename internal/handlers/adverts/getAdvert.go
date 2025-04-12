package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}

	advertID := c.Params("id")

	advert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	if advert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}
	response := fiber.Map{
		"status": "success",
		"data":   advert,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
