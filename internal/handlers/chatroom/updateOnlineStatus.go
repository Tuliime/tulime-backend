package chatroom

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateOnlineStatus = func(c *fiber.Ctx) error {
	onlineStatus := models.OnlineStatus{}

	if err := c.BodyParser(&onlineStatus); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if onlineStatus.UserID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID!")
	}

	savedOnlineStatus, err := onlineStatus.FindByUser(onlineStatus.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedOnlineStatus.ID == "" {
		onlineStatus, err = onlineStatus.Create(onlineStatus)
		if err != nil {
			log.Println("err.Error(): ", err.Error())
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	} else {
		onlineStatus, err = savedOnlineStatus.Update()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	log.Printf("published event onlineStatus : %+v", onlineStatus)
	events.EB.Publish("onlineStatus", onlineStatus)

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    onlineStatus,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
