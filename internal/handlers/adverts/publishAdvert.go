package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PublishAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	advertID := c.Params("id")

	savedAdvert, err := advert.Find(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	if savedAdvert.IsPublished {
		return fiber.NewError(fiber.StatusNotFound, "Advert is already published!")
	}

	savedAdvert.IsPublished = true
	publishedAdvert, err := savedAdvert.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	updatedAdvert, err := advert.FindOne(publishedAdvert.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert published successfully!",
		"data":    updatedAdvert,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
