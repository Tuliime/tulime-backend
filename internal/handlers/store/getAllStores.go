package store

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllStores = func(c *fiber.Ctx) error {
	store := models.Store{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	inCludeCursorParam := c.Query("includeCursor", "false")
	direction := c.Query("direction")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	includeCursor, err := strconv.ParseBool(inCludeCursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if direction == "" {
		direction = "BACKWARD"
	}

	stores, err := store.FindAll(limit+1, cursorParam, includeCursor, direction)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor, nextCursor string
	var hasPrevItems, hasNextItems bool

	if len(stores) > 0 && direction == "BACKWARD" {
		if len(stores) > int(limit) {
			stores = stores[:len(stores)-1] // Remove last element
			prevCursor = stores[0].ID
			hasPrevItems = true
			if cursorParam != "" {
				nextCursor = stores[len(stores)-1].ID
				hasNextItems = true
			}
		} else {
			prevCursor = ""
			hasPrevItems = false
			if cursorParam != "" {
				nextCursor = stores[len(stores)-1].ID
				hasNextItems = true
			}
		}
	}

	if len(stores) > 0 && direction == "FORWARD" {
		if len(stores) > int(limit) {
			stores = stores[:len(stores)-1] // Remove last element
			nextCursor = stores[len(stores)-1].ID
			hasNextItems = true
			if cursorParam != "" {
				prevCursor = stores[0].ID
				hasPrevItems = true
			}
		} else {
			nextCursor = ""
			hasNextItems = false
			if cursorParam != "" {
				prevCursor = stores[0].ID
				hasPrevItems = true
			}
		}
	}

	pagination := map[string]interface{}{
		"limit":         limit,
		"prevCursor":    prevCursor,
		"nextCursor":    nextCursor,
		"includeCursor": includeCursor,
		"hasNextItems":  hasNextItems,
		"hasPrevItems":  hasPrevItems,
		"direction":     direction,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       stores,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
