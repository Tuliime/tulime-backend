package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateAdvertInventory = func(c *fiber.Ctx) error {
	advertInventoryID := c.Params("advertInventoryID")
	advertInventory := models.AdvertInventory{}

	savedAdvertInventory, err := advertInventory.FindOne(advertInventoryID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvertInventory.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert inventory of provided id doesn't exist!")
	}

	if err := c.BodyParser(&advertInventory); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advertInventory.Quantity == 0 || advertInventory.Unit == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing Quantity/Unit!")
	}

	savedAdvertInventory.Quantity = advertInventory.Quantity
	savedAdvertInventory.Unit = advertInventory.Unit

	updatedAdvertInventory, err := savedAdvertInventory.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert inventory updated successfully!",
		"data":    updatedAdvertInventory,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
