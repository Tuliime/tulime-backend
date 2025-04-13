package store

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetStore = func(c *fiber.Ctx) error {
	store := models.Store{}
	advert := models.Advert{}
	storeID := c.Params("id")
	includeAdvertStr := c.Query("includeAdverts", "false")

	store, err := store.FindOne(storeID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	includeAdverts, err := strconv.ParseBool(includeAdvertStr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var adverts []models.Advert
	if includeAdverts {
		adverts, err = advert.FindByStore(storeID, 20, "")
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
