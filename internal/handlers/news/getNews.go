package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetNews = func(c *fiber.Ctx) error {
	news := models.News{}
	newsID := c.Params("id")

	news, err := news.FindOne(newsID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   news,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
