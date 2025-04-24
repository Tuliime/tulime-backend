package messenger

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type MessengerIDs struct {
	MessengerIDList string `json:"messengerIDList"`
}

// TODO: consider sending updated messages via live sse
// TODO: consider saving unread messages in batches
var UpdateMessagesAsRead = func(c *fiber.Ctx) error {
	messenger := models.Messenger{}
	messageIDs := MessengerIDs{}
	// messengerRoomID := c.Params("messengerRoomID")
	userID := c.Locals("userID")

	// if messengerRoomID == "" {
	// 	return fiber.NewError(fiber.StatusBadRequest, "Please provide messengerRoomID!")
	// }

	if err := c.BodyParser(&messageIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var messengerIDList []string
	if messageIDs.MessengerIDList != "" {
		err := json.Unmarshal([]byte(messageIDs.MessengerIDList), &messengerIDList)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest,
				"Invalid messengerIDList format! Must be a JSON stringified array of strings.")
		}
	}
	fmt.Printf("messengerIDList: %v\n", messengerIDList)

	for _, messengerID := range messengerIDList {
		savedMessage, err := messenger.Find(messengerID)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}

		if savedMessage.ID == "" {
			log.Println("Message of provided id does'nt exist!")
			continue
		}
		if savedMessage.IsRead {
			log.Println("Message is already updated as read!")
			continue
		}
		if savedMessage.SenderID == userID {
			log.Println("User is sender of the message!")
			continue
		}
		savedMessage.IsRead = true

		updateMessage, err := savedMessage.Update()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		log.Printf("updated message:  %v", updateMessage)
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Messenger updated successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
