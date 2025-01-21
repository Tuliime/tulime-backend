package chatroom

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetChat = func(c *fiber.Ctx) error {
	chatroom := models.Chatroom{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	inCludeCursorParam := c.Query("includeCursor", "false")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}

	includeCursor, err := strconv.ParseBool(inCludeCursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	chatMessages, err := chatroom.FindAll(limit, cursorParam, includeCursor)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var repliedMessages []models.Chatroom

	for _, chatMessage := range chatMessages {
		if chatMessage.Reply == "" {
			continue
		}
		reply, err := chatroom.FindReply(chatMessage.Reply)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		repliedMessages = append(repliedMessages, reply)
	}

	var prevCursor string
	if len(chatMessages) > 0 {
		prevCursor = chatMessages[0].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	chatMap := fiber.Map{
		"chats":   chatMessages,
		"replies": repliedMessages,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       chatMap,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
