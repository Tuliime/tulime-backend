package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PostAdvertInventory = func(c *fiber.Ctx) error {
	advertID := c.Params("id")
	advert := models.Advert{}
	advertInventory := models.AdvertInventory{}

	savedAdvert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	savedAdvertInventory, err := advertInventory.FindByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvertInventory.ID != "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id already has inventory!")
	}

	advertInventory.AdvertID = advertID

	if err := c.BodyParser(&advertInventory); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advertInventory.Quantity == 0 || advertInventory.Unit == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Quantity/Unit!")
	}

	newAdvertInventory, err := advertInventory.Create(advertInventory)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert inventory created successfully!",
		"data":    newAdvertInventory,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
