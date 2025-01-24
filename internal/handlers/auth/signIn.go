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

	accessToken, err := packages.SignJWTToken(user.ID, "accessToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := packages.SignJWTToken(user.ID, "refreshToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	session := models.Session{UserID: user.ID, AccessToken: accessToken,
		RefreshToken: refreshToken, GeneratedVia: "sign in"}
	if _, err := session.Create(session); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userMap := map[string]interface{}{
		"id":        user.ID,
		"name":      user.Name,
		"telNumber": user.TelNumber,
		"role":      user.Role,
		"imageUrl":  user.ImageUrl,
	}
	response := map[string]interface{}{
		"status":       "success",
		"message":      "Sign in successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
