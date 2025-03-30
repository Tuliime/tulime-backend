package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAdvertsByUser = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	userID := c.Params("userID")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	adverts, err := advert.FindByUser(userID, limit+1, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	var hasPrevItems bool

	if len(adverts) > 0 && len(adverts) > int(limit) {
		adverts = adverts[:len(adverts)-1] // Remove last element
		prevCursor = adverts[0].ID
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
