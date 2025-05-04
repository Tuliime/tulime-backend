package auth

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/handlers/location"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var SignIn = func(c *fiber.Ctx) error {
	user := models.User{}
	device := c.Get("X-Device")
	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid client type!")
	}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	password := user.Password

	hasEmail := user.Email != ""
	hasTelNumber := user.TelNumber != 0

	if !hasEmail && !hasTelNumber {
		return fiber.NewError(fiber.StatusBadRequest, "Email or telephone number is required")
	}

	if password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Password is required")
	}

	if hasTelNumber {
		userByNumber, err := user.FindByTelNumber(user.TelNumber)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if userByNumber.ID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid telephone number/password!")
		}
		user = userByNumber
	}

	if hasEmail {
		userByEmail, err := user.FindByEmail(user.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if userByEmail.ID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid email/password!")
		}
		user = userByEmail
	}

	passwordMatches, err := user.PasswordMatches(password)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if !passwordMatches {
		if hasTelNumber {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid telephone number/password!")
		} else if hasEmail {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid email/password!")
		}
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
		RefreshToken: refreshToken, GeneratedVia: "sign in", Device: device,
		LocationID: location.ID}
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
		"message":      "Sign in successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
