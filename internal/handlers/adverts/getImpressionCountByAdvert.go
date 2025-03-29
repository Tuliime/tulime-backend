package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetImpressionCountByAdvert = func(c *fiber.Ctx) error {
	advertImpression := models.AdvertImpression{}

	advertID := c.Params("id")

	advertViewCount, err := advertImpression.FindCountByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	advertImpressionCountMap := fiber.Map{
		"count": advertViewCount,
	}

	response := fiber.Map{
		"status": "success",
		"data":   advertImpressionCountMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
