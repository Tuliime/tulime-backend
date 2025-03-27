package store

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateStore = func(c *fiber.Ctx) error {
	store := models.Store{}
	storeID := c.Params("id")

	if err := c.BodyParser(&store); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if store.Name == "" || store.Description == "" || store.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/userID/name/description!")
	}

	savedStore, err := store.FindOne(storeID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedStore.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Store of provided id is not found!")
	}

	// Update changed fields
	if store.Name != "" {
		savedStore.Name = store.Name
	}
	if store.Description != "" {
		savedStore.Description = store.Description
	}
	if store.Website != "" {
		savedStore.Website = store.Website
	}
	if store.Email != "" {
		savedStore.Email = store.Email
	}
	if store.Type != "" {
		savedStore.Type = store.Type
	}

	updatedStore, err := savedStore.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Store updated successfully!",
		"data":    updatedStore,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
