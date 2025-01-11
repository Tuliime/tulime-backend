package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var SignIn = func(c *fiber.Ctx) error {
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	password := user.Password

	if user.TelNumber == 0 || user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing username/telNumber/password!")
	}

	user, err := user.FindByTelNumber(user.TelNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid telNumber/password!")
	}

	passwordMatches, err := user.PasswordMatches(password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !passwordMatches {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid telNumber/password!")
	}

	accessToken, err := packages.SignJWTToken(user.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userMap := map[string]interface{}{
		"id":        user.ID,
		"name":      user.Name,
		"telNumber": user.TelNumber,
		"role":      user.Role,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Sign in successfully",
		"accessToken": accessToken,
		"user":        userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
