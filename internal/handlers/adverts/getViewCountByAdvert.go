package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetViewCountByAdvert = func(c *fiber.Ctx) error {
	advertView := models.AdvertView{}

	advertID := c.Params("id")

	advertViewCount, err := advertView.FindCountByAdvert(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	advertViewCountMap := fiber.Map{
		"count": advertViewCount,
	}

	response := fiber.Map{
		"status": "success",
		"data":   advertViewCountMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
