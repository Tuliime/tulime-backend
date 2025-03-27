package adverts

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllAdverts = func(c *fiber.Ctx) error {
	advert := models.Advert{}
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

	adverts, err := advert.FindAll(limit+1, cursorParam, includeCursor, direction)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor, nextCursor string
	var hasPrevItems, hasNextItems bool

	if len(adverts) > 0 && direction == "BACKWARD" {
		if len(adverts) > int(limit) {
			adverts = adverts[:len(adverts)-1] // Remove last element
			prevCursor = adverts[0].ID
			hasPrevItems = true
			if cursorParam != "" {
				nextCursor = adverts[len(adverts)-1].ID
				hasNextItems = true
			}
		} else {
			prevCursor = ""
			hasPrevItems = false
			if cursorParam != "" {
				nextCursor = adverts[len(adverts)-1].ID
				hasNextItems = true
			}
		}
	}

	if len(adverts) > 0 && direction == "FORWARD" {
		if len(adverts) > int(limit) {
			adverts = adverts[:len(adverts)-1] // Remove last element
			nextCursor = adverts[len(adverts)-1].ID
			hasNextItems = true
			if cursorParam != "" {
				prevCursor = adverts[0].ID
				hasPrevItems = true
			}
		} else {
			nextCursor = ""
			hasNextItems = false
			if cursorParam != "" {
				prevCursor = adverts[0].ID
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

	storeMap := fiber.Map{
		"stores": adverts,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       storeMap,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
