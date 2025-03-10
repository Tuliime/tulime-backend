package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllNews = func(c *fiber.Ctx) error {
	news := models.News{}
	limitParam := c.Query("limit")
	categoryParam := c.Query("category")
	cursorParam := c.Query("cursor")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}
	if categoryParam == "" {
		categoryParam = ""
	}

	allNews, err := news.FindAll(limit, categoryParam, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	if len(allNews) > 0 {
		prevCursor = allNews[len(allNews)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allNews,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
