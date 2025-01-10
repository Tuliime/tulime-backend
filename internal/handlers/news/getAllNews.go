package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllNews = func(c *fiber.Ctx) error {
	news := models.News{}
	limitParam := c.Query("limit")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	allNews, err := news.FindAll(limit)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	pagination := map[string]interface{}{
		"limit": limit,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allNews,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
