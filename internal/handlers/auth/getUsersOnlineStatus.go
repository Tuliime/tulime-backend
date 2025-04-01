package auth

import (
	"encoding/base64"
	"encoding/json"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetUsersOnlineStatus = func(c *fiber.Ctx) error {

	onlineStatus := models.OnlineStatus{}
	userIDListEncoding := c.Query("userIDListEncoding")

	decodedBytes, err := base64.StdEncoding.DecodeString(userIDListEncoding)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userIDListStr := string(decodedBytes)
	if userIDListStr == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "Provided encoding without content!")
	}

	var userIDList []string
	if userIDListStr != "" {
		err = json.Unmarshal(decodedBytes, &userIDList)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid format! encoded data must array of strings.")
		}
	}

	statuses, err := onlineStatus.FindByUserList(userIDList)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   statuses,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
