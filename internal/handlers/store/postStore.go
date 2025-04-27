package store

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// Allow one store per user
var PostStore = func(c *fiber.Ctx) error {
	store := models.Store{}
	if err := c.BodyParser(&store); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid string type for userID")
	}

	if store.Name == "" || store.Description == "" || store.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/userID/name/description!")
	}

	savedStore, err := store.FindByUser(userID, 1, "")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(savedStore) > 0 {
		return fiber.NewError(fiber.StatusBadRequest, "You already created a store!")
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
