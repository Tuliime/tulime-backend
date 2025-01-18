package chatbot

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetChatByUser = func(c *fiber.Ctx) error {
	chatbot := models.Chatbot{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	userID := c.Params("userID")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}

	chats, err := chatbot.FindByUser(userID, limit, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	if len(chats) > 0 {
		prevCursor = chats[0].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       chats,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
