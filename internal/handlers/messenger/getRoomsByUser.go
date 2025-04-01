package messenger

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

// Gets last message for each room by userID
// And so it is essentially to the user rooms
var GetRoomsByUser = func(c *fiber.Ctx) error {
	messenger := models.Messenger{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	userID := c.Params("userID")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	messengerRooms, err := messenger.FindRoomsByUser(userID, limit+1, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	var hasPrevItems bool

	if len(messengerRooms) > 0 && len(messengerRooms) > int(limit) {
		messengerRooms = messengerRooms[:len(messengerRooms)-1] // Remove last element
		prevCursor = messengerRooms[0].ID
		hasPrevItems = true
	} else {
		prevCursor = ""
		hasPrevItems = false
	}

	pagination := fiber.Map{
		"limit":        limit,
		"prevCursor":   prevCursor,
		"hasPrevItems": hasPrevItems,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       organizeRoomResponse(messengerRooms),
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func organizeRoomResponse(messengerRooms []models.Messenger) []Message {
	var response []Message

	for _, room := range messengerRooms {
		response = append(response, Message{
			ID:              room.ID,
			MessengerRoomID: room.MessengerRoomID,
			SenderID:        room.SenderID,
			RecipientID:     room.RecipientID,
			Text:            room.Text,
			Reply:           room.Reply,
			File: models.MessengerFile{
				ID:          room.File.ID,
				MessengerID: room.File.MessengerID,
				URL:         room.File.URL,
				Path:        room.File.Path,
				CreatedAt:   room.File.CreatedAt,
				UpdatedAt:   room.File.UpdatedAt,
			},
			Tag:       room.Tag,
			IsRead:    room.IsRead,
			SentAt:    room.SentAt,
			ArrivedAt: room.ArrivedAt,
			CreatedAt: room.CreatedAt,
			UpdatedAt: room.UpdatedAt,
			Sender: User{
				ID:             room.Sender.ID,
				Name:           room.Sender.Name,
				TelNumber:      room.Sender.TelNumber,
				Role:           room.Sender.Role,
				ImageUrl:       room.Sender.ImageUrl,
				ImagePath:      room.Sender.ImagePath,
				ProfileBgColor: room.Sender.ProfileBgColor,
				ChatroomColor:  room.Sender.ChatroomColor,
				CreatedAt:      room.Sender.CreatedAt,
				UpdatedAt:      room.Sender.UpdatedAt,
			},
			Recipient: User{
				ID:             room.Recipient.ID,
				Name:           room.Recipient.Name,
				TelNumber:      room.Recipient.TelNumber,
				Role:           room.Recipient.Role,
				ImageUrl:       room.Recipient.ImageUrl,
				ImagePath:      room.Recipient.ImagePath,
				ProfileBgColor: room.Recipient.ProfileBgColor,
				ChatroomColor:  room.Recipient.ChatroomColor,
				CreatedAt:      room.Recipient.CreatedAt,
				UpdatedAt:      room.Recipient.UpdatedAt,
			},
		})
	}

	return response
}
