package chatroom

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetChat = func(c *fiber.Ctx) error {
	chatroom := models.Chatroom{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	var cursor string

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursor = ""
	}

	chatMessages, err := chatroom.FindAll(limit, cursor)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// TODO: To add a cursor in the pagination
	pagination := map[string]interface{}{
		"limit": limit,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       chatMessages,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
