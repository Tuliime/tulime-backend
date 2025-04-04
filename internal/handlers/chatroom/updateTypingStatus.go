package chatroom

import (
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateTypingStatus = func(c *fiber.Ctx) error {
	typingStatus := TypingStatus{}
	user := models.User{}
	typingStatusInput := struct {
		UserID          string `json:"userID"`
		StartedTypingAt string `json:"startedTypingAt"`
		RecipientID     string `json:"recipientID"`
		Type            string `json:"type"`
	}{}

	if err := c.BodyParser(&typingStatusInput); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if typingStatusInput.UserID == "" || typingStatusInput.RecipientID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing userID/recipientID!")
	}

	user, err := user.FindOne(typingStatusInput.UserID)
	if typingStatusInput.UserID == "" {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	typingStatus.Type = typingStatusInput.Type
	typingStatus.UserID = typingStatusInput.UserID
	typingStatus.RecipientID = typingStatusInput.RecipientID
	typingStatus.User = User{ID: user.ID,
		Name:           user.Name,
		TelNumber:      user.TelNumber,
		Role:           user.Role,
		ImageUrl:       user.ImageUrl,
		ImagePath:      user.ImagePath,
		ProfileBgColor: user.ProfileBgColor,
		ChatroomColor:  user.ChatroomColor,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt}
	typingStatus.StartedTypingAt = time.Now()

	events.EB.Publish("typingStatus", typingStatus)

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    typingStatus,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
