package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

type Passwords struct {
	CurrentPassword string `json:"currentPassword"`
	NewPassword     string `json:"newPassword"`
}

var ChangePassword = func(c *fiber.Ctx) error {
	user := models.User{ID: c.Params("id")}
	passwords := Passwords{}

	if err := c.BodyParser(&passwords); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if passwords.CurrentPassword == "" || passwords.NewPassword == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing current/new password!")
	}

	if passwords.CurrentPassword == passwords.NewPassword {
		return fiber.NewError(fiber.StatusBadRequest, "New password is same as current password")
	}

	user, err := user.FindOne(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "We couldn't find user of provided id!")
	}

	currentPasswordMatchesSavedOne, err := user.PasswordMatches(passwords.CurrentPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !currentPasswordMatchesSavedOne {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid current password!")
	}

	hashedPassword, err := user.HashPassword(passwords.NewPassword)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user.Password = hashedPassword

	if err := user.Update(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "Password changed successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
