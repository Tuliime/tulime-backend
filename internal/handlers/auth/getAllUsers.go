package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetAllUsers = func(c *fiber.Ctx) error {
	user := models.User{}
	users, err := user.FindAll()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	filteredUSers := make([]map[string]interface{}, len(users))
	for i, user := range users {
		filteredUSers[i] = map[string]interface{}{
			"id":             user.ID,
			"name":           user.Name,
			"role":           user.Role,
			"telNumber":      user.TelNumber,
			"imageUrl":       user.ImageUrl,
			"profileBgColor": user.ProfileBgColor,
			"chatroomColor":  user.ChatroomColor,
			"createdAt":      user.CreatedAt,
			"updatedAt":      user.UpdatedAt,
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   filteredUSers,
	})
}
