package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	advertID := c.Params("id")
	if err := c.BodyParser(&advert); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advert.ProductName == "" || advert.ProductDescription == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"Missing ProductName/ProductDescription!")
	}

	savedAdvert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	savedAdvert.ProductName = advert.ProductName
	savedAdvert.ProductDescription = advert.ProductDescription

	updatedAdvert, err := savedAdvert.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert updated successfully!",
		"data":    updatedAdvert,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
