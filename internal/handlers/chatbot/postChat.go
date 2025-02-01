package chatbot

import (
	"fmt"
	"log"
	"time"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type PostChatInput struct {
	ID        string `validate:"string"`
	Message   string `validate:"string"`
	WrittenBy string `validate:"string"`
	PostedAt  string `validate:"string"`
}

var PostChat = func(c *fiber.Ctx) error {
	chatbot := models.Chatbot{UserID: c.Params("userID")}

	if err := c.BodyParser(&chatbot); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input PostChatInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	userPostedAt, err := time.Parse(time.RFC3339, input.PostedAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sentAt format! Must be an ISO 8601 string.")
	}
	fmt.Printf("parsedPostedAt: %v\n", userPostedAt)

	newUserChat, err := chatbot.Create(models.Chatbot{ID: input.ID, UserID: c.Params("userID"),
		Message: input.Message, WrittenBy: input.WrittenBy, PostedAt: userPostedAt})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// TODO: To make an api call to openai API here

	botPostedAt := time.Now()
	botMessage := "This is the apparent ai response from openai"

	newBotChat, err := chatbot.Create(models.Chatbot{UserID: c.Params("userID"),
		Message: botMessage, WrittenBy: "bot", PostedAt: botPostedAt})
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Created successfully!",
		"data": fiber.Map{
			"user": newUserChat,
			"bot":  newBotChat,
		},
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
