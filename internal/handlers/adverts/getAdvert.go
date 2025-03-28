package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	advertView := models.AdvertView{}
	advertImpression := models.AdvertImpression{}

	advertID := c.Params("id")
	userID := c.Locals("userID")

	advert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if advert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	var advertViewCount, advertImpressionCount int64
	if advert.UserID == userID {
		advertViewCount, err = advertView.FindCountByAdvert(advertID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		advertImpressionCount, err = advertImpression.FindCountByAdvert(advertID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	advertMap := fiber.Map{
		"advert": advert,
		"count": fiber.Map{
			"viewCount":       advertViewCount,
			"impressionCount": advertImpressionCount,
		},
	}

	response := fiber.Map{
		"status": "success",
		"data":   advertMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
