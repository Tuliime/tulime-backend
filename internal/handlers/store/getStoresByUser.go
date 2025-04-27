package store

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetStoresByUser = func(c *fiber.Ctx) error {
	store := models.Store{}
	cursorParam := c.Query("cursor")
	userID := c.Params("userID")
	includeAdvertStr := c.Query("includeAdverts", "false")

	stores, err := store.FindByUser(userID, 1, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if len(stores) > 0 {
		store = stores[0]
	}

	includeAdverts, err := strconv.ParseBool(includeAdvertStr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	advert := models.Advert{}
	var adverts []models.Advert

	if includeAdverts && len(stores) > 0 {
		adverts, err = advert.FindByStore(stores[0].ID, 20, "")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		for i := range adverts {
			store.Advert = append(store.Advert, &adverts[i])
		}
	}

	response := fiber.Map{
		"status": "success",
		"data":   store,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
