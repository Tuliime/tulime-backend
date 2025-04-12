package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetAdvertAnalytics = func(c *fiber.Ctx) error {
	advertView := models.AdvertView{}
	advertImpression := models.AdvertImpression{}
	advertID := c.Params("id")

	var advertViewCount, advertImpressionCount int64
	advertViewCount, err := advertView.FindCountByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	advertImpressionCount, err = advertImpression.FindCountByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	advertAnalytics := fiber.Map{
		"advertID":        advertID,
		"viewCount":       advertViewCount,
		"impressionCount": advertImpressionCount,
	}

	response := fiber.Map{
		"status": "success",
		"data":   advertAnalytics,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
