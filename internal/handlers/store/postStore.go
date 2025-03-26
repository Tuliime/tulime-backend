package store

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PostStore = func(c *fiber.Ctx) error {
	store := models.Store{}
	if err := c.BodyParser(&store); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if store.Name == "" || store.Description == "" || store.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/userID/name/description!")
	}

	newStore, err := store.Create(store)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Store created successfully!",
		"data":    newStore,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
