package store

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetStoresByUser = func(c *fiber.Ctx) error {
	store := models.Store{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	userID := c.Params("userID")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	stores, err := store.FindByUSer(userID, limit+1, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	var hasPrevItems bool

	if len(stores) > 0 && len(stores) > int(limit) {
		stores = stores[:len(stores)-1] // Remove last element
		prevCursor = stores[0].ID
		hasPrevItems = true
	} else {
		prevCursor = ""
		hasPrevItems = false
	}

	pagination := fiber.Map{
		"limit":        limit,
		"prevCursor":   prevCursor,
		"hasPrevItems": hasPrevItems,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       stores,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
