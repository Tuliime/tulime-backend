package chatbot

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var DeleteChat = func(c *fiber.Ctx) error {
	chat := models.Chatbot{}
	chatId := c.Params("id")

	savedChat, err := chat.FindOne(chatId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedChat.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Chat of provided id is not found!")
	}

	err = chat.Delete(chatId)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Chat deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
