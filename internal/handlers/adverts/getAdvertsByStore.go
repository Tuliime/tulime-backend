package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAdvertsByStore = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	storeID := c.Params("storeID")

	limitParam := c.Query("limit")
	cursor := c.Query("cursor")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	adverts, err := advert.FindByStore(storeID, limit+1, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	response := fiber.Map{
		"status":     "success",
		"data":       adverts,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
