package chatbot

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type PostChatInput struct {
	Prompt string `validate:"string"`
}

var PostChat = func(c *fiber.Ctx) error {
	chatbot := models.Chatbot{UserID: c.Params("userId")}

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

	// TODO: To make an api call to openai API here

	chatbot.AIResponse = "This is the apparent ai response from openai"

	newChat, err := chatbot.Create(chatbot)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Created successfully!",
		"data":    newChat,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
