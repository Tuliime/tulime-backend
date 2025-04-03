package auth

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/handlers/location"
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

	device := c.Get("X-Device")
	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid client type!")
	}

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

	location, err := location.GetUserLocationByIP(user.ID, clientIP)
	if err != nil {
		log.Printf("Error getting location ID:  %+v", err)
	}

	session.UserID = user.ID
	session.AccessToken = accessToken
	session.RefreshToken = sigInInput.RefreshToken
	session.GeneratedVia = "sign in with refresh token"
	session.Device = device
	session.LocationID = location.ID

	if _, err := session.Create(session); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userMap := map[string]interface{}{
		"id":             user.ID,
		"name":           user.Name,
		"telNumber":      user.TelNumber,
		"role":           user.Role,
		"imageUrl":       user.ImageUrl,
		"profileBgColor": user.ProfileBgColor,
		"chatroomColor":  user.ChatroomColor,
		"createdAt":      user.CreatedAt,
		"updatedAt":      user.UpdatedAt,
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
