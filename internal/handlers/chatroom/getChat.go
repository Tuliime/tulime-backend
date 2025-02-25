package chatroom

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetChat = func(c *fiber.Ctx) error {
	chatroom := models.Chatroom{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	inCludeCursorParam := c.Query("includeCursor", "false")
	direction := c.Query("direction")

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

	if direction == "" {
		direction = "BACKWARD"
	}

	chatMessages, err := chatroom.FindAll(limit, cursorParam, includeCursor, direction)
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
			if err.Error() == constants.RECORD_NOT_FOUND_ERROR {
				continue
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		repliedMessages = append(repliedMessages, reply)
	}

	var prevCursor, nextCursor string
	if len(chatMessages) > 0 && direction == "BACKWARD" {
		prevCursor = chatMessages[0].ID
	}
	if len(chatMessages) > 0 && direction == "FORWARD" {
		nextCursor = chatMessages[len(chatMessages)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":         limit,
		"prevCursor":    prevCursor,
		"nextCursor":    nextCursor,
		"includeCursor": includeCursor,
	}

	chatMap := fiber.Map{
		"messages": chatMessages,
		"replies":  repliedMessages,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       chatMap,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
