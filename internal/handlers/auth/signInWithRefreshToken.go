package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type SignInWithRefreshTokenInput struct {
	UserID       string `json:"userID"`
	RefreshToken string `json:"refreshToken"`
}

var SignInWithRefreshToken = func(c *fiber.Ctx) error {
	user := models.User{}
	session := models.Session{}
	sigInInput := SignInWithRefreshTokenInput{}

	if err := c.BodyParser(&sigInInput); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if sigInInput.UserID == "" || sigInInput.RefreshToken == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing userID/refreshToken!")
	}

	user, err := user.FindOne(sigInInput.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "User of provided id does not exist!")
	}

	userIDFromToken, err := packages.DecodeJWT(sigInInput.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusNotFound, err.Error())
	}

	savedRefreshToken, err := session.FindByRefreshToken(sigInInput.RefreshToken)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedRefreshToken.ID == "" || savedRefreshToken.UserID != user.ID || savedRefreshToken.UserID != userIDFromToken {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid refresh token!")
	}

	accessToken, err := packages.SignJWTToken(user.ID, "accessToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	session.UserID = user.ID
	session.AccessToken = accessToken
	session.RefreshToken = sigInInput.RefreshToken
	session.GeneratedVia = "sign in with refresh token"

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
		"message":      "Auto sign in successful",
		"accessToken":  accessToken,
		"refreshToken": sigInInput.RefreshToken,
		"user":         userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
