package auth

import (
	"log"
	"time"

	"github.com/Tuliime/tulime-backend/internal/handlers/location"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var ResetPassword = func(c *fiber.Ctx) error {
	user := models.User{}
	otp := models.OTP{}

	device := c.Get("X-Device")
	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid client type!")
	}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newPassword := user.Password
	OPT := c.Params("otp")

	if user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide your new password!")
	}

	otp, err := otp.FindByOTP(OPT)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if otp.IsUsed || !otp.IsVerified || otp.ExpiresAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "OPT is already used or unverified or expired!")
	}

	user, err = user.FindOne(otp.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if err := user.ResetPassword(newPassword); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	otp.ExpiresAt = time.Now()
	otp.IsUsed = true
	if _, err := otp.Update(); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessToken, err := packages.SignJWTToken(user.ID, "accessToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := packages.SignJWTToken(user.ID, "refreshToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	location, err := location.GetUserLocationByIP(user.ID, clientIP)
	if err != nil {
		log.Printf("Error getting location ID:  %+v", err)
	}

	session := models.Session{UserID: user.ID, AccessToken: accessToken,
		RefreshToken: refreshToken, GeneratedVia: "reset password",
		Device: device, LocationID: location.ID}
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
		"message":      "Password reset successfully!",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
