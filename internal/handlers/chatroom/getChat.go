package chatroom

import (
	"encoding/json"
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

	chatMessages, err := chatroom.FindAll(limit+1, cursorParam, includeCursor, direction)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var repliedMessages []models.Chatroom
	var messages []Message

	for _, chatMessage := range chatMessages {
		var repliedMessage any = nil
		var repliedMessageFile any = nil

		if chatMessage.Reply != "" {
			reply, err := chatroom.FindReply(chatMessage.Reply)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			if reply.File.ID != "" {
				var dimensions models.ImageDimensions
				if err := json.Unmarshal(reply.File.Dimensions, &dimensions); err != nil {
					return fiber.NewError(fiber.StatusInternalServerError, err.Error())
				}
				repliedMessageFile = File{
					ID:         reply.File.ID,
					ChatroomID: reply.File.ChatroomID,
					URL:        reply.File.URL,
					Path:       reply.File.Path,
					Dimensions: dimensions,
					CreatedAt:  reply.File.CreatedAt,
					UpdatedAt:  reply.File.UpdatedAt,
					DeletedAt:  reply.File.DeletedAt,
				}
			}
			repliedMessage = Message{
				ID:             reply.ID,
				UserID:         reply.UserID,
				Text:           reply.Text,
				Reply:          reply.Reply,
				RepliedMessage: nil,
				File:           repliedMessageFile,
				Mention:        reply.Mention,
				SentAt:         reply.SentAt,
				ArrivedAt:      reply.ArrivedAt,
				CreatedAt:      reply.CreatedAt,
				UpdatedAt:      reply.UpdatedAt,
				DeletedAt:      reply.DeletedAt,
				User: User{
					ID:             reply.User.ID,
					Name:           reply.User.Name,
					TelNumber:      reply.User.TelNumber,
					Role:           reply.User.Role,
					ImageUrl:       reply.User.ImageUrl,
					ImagePath:      reply.User.ImagePath,
					ProfileBgColor: reply.User.ProfileBgColor,
					ChatroomColor:  reply.User.ChatroomColor,
					CreatedAt:      reply.User.CreatedAt,
					UpdatedAt:      reply.User.UpdatedAt,
				},
			}
			repliedMessages = append(repliedMessages, reply)
		}

		var chatMessageFile any = nil
		if chatMessage.File.ID != "" {
			var dimensions models.ImageDimensions
			if err := json.Unmarshal(chatMessage.File.Dimensions, &dimensions); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			chatMessageFile = File{
				ID:         chatMessage.File.ID,
				ChatroomID: chatMessage.File.ChatroomID,
				URL:        chatMessage.File.URL,
				Path:       chatMessage.File.Path,
				Dimensions: dimensions,
				CreatedAt:  chatMessage.File.CreatedAt,
				UpdatedAt:  chatMessage.File.UpdatedAt,
				DeletedAt:  chatMessage.File.DeletedAt,
			}
		}
		messages = append(messages, Message{
			ID:             chatMessage.ID,
			UserID:         chatMessage.UserID,
			Text:           chatMessage.Text,
			Reply:          chatMessage.Reply,
			RepliedMessage: repliedMessage,
			File:           chatMessageFile,
			Mention:        chatMessage.Mention,
			SentAt:         chatMessage.SentAt,
			ArrivedAt:      chatMessage.ArrivedAt,
			CreatedAt:      chatMessage.CreatedAt,
			UpdatedAt:      chatMessage.UpdatedAt,
			DeletedAt:      chatMessage.DeletedAt,
			User: User{
				ID:             chatMessage.User.ID,
				Name:           chatMessage.User.Name,
				TelNumber:      chatMessage.User.TelNumber,
				Role:           chatMessage.User.Role,
				ImageUrl:       chatMessage.User.ImageUrl,
				ImagePath:      chatMessage.User.ImagePath,
				ProfileBgColor: chatMessage.User.ProfileBgColor,
				ChatroomColor:  chatMessage.User.ChatroomColor,
				CreatedAt:      chatMessage.User.CreatedAt,
				UpdatedAt:      chatMessage.User.UpdatedAt,
			},
		})
	}

	var prevCursor, nextCursor string
	var hasPrevItems, hasNextItems bool

	if len(messages) > 0 && direction == "BACKWARD" {
		if len(messages) > int(limit) {
			messages = messages[1:] // Remove first element
			prevCursor = messages[0].ID
			hasPrevItems = true
			if cursorParam != "" {
				nextCursor = messages[len(messages)-1].ID
				hasNextItems = true
			}
		} else {
			prevCursor = ""
			hasPrevItems = false
			if cursorParam != "" {
				nextCursor = messages[len(messages)-1].ID
				hasNextItems = true
			}
		}
	}

	if len(messages) > 0 && direction == "FORWARD" {
		if len(messages) > int(limit) {
			messages = messages[:len(messages)-1] // Remove last element
			nextCursor = messages[len(messages)-1].ID
			hasNextItems = true
			if cursorParam != "" {
				prevCursor = messages[0].ID
				hasPrevItems = true
			}
		} else {
			nextCursor = ""
			hasNextItems = false
			if cursorParam != "" {
				prevCursor = messages[0].ID
				hasPrevItems = true
			}
		}
	}

	pagination := map[string]interface{}{
		"limit":         limit,
		"prevCursor":    prevCursor,
		"nextCursor":    nextCursor,
		"includeCursor": includeCursor,
		"hasNextItems":  hasNextItems,
		"hasPrevItems":  hasPrevItems,
		"direction":     direction,
	}

	chatMap := fiber.Map{
		"messages": messages,
		"replies":  repliedMessages,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       chatMap,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
