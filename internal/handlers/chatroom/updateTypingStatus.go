package chatroom

import (
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/gofiber/fiber/v2"
)

var UpdateTypingStatus = func(c *fiber.Ctx) error {
	typingStatus := TypingStatus{}
	typingStatusInput := struct {
		UserID          string `json:"userID"`
		StartedTypingAt string `json:"startedTypingAt"`
	}{}

	if err := c.BodyParser(&typingStatusInput); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if typingStatusInput.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID!")
	}

	typingStatus.UserID = typingStatusInput.UserID
	typingStatus.StartedTypingAt = time.Now()

	events.EB.Publish("typingStatus", typingStatus)

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    typingStatus,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
