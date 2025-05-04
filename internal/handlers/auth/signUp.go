package auth

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/handlers/location"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var SignUp = func(c *fiber.Ctx) error {
	user := models.User{}
	device := c.Get("X-Device")
	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid client type!")
	}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.Name == "" || user.TelNumber == 0 || user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing username/telephone number/password!")
	}

	userByNumber, err := user.FindByTelNumber(user.TelNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if userByNumber.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Telephone number already registered!")
	}

	userByEmail, err := user.FindByEmail(user.Email)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if userByEmail.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email is already registered!")
	}

	err = user.SetRole("user")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newUser, err := user.Create(user)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	log.Println("newUser: ", newUser)

	accessToken, err := packages.SignJWTToken(newUser.ID, "accessToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	refreshToken, err := packages.SignJWTToken(newUser.ID, "refreshToken")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	location, err := location.GetUserLocationByIP(newUser.ID, clientIP)
	if err != nil {
		log.Printf("Error getting location ID:  %+v", err)
	}

	session := models.Session{UserID: newUser.ID, AccessToken: accessToken,
		RefreshToken: refreshToken, GeneratedVia: "sign up", Device: device,
		LocationID: location.ID}
	if _, err := session.Create(session); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newUserResponseData := map[string]interface{}{
		"id":             newUser.ID,
		"name":           newUser.Name,
		"telNumber":      newUser.TelNumber,
		"role":           newUser.Role,
		"imageUrl":       newUser.ImageUrl,
		"profileBgColor": newUser.ProfileBgColor,
		"chatroomColor":  user.ChatroomColor,
		"createdAt":      newUser.CreatedAt,
		"updatedAt":      newUser.UpdatedAt,
	}
	response := map[string]interface{}{
		"status":       "success",
		"message":      "Account created successfully",
		"accessToken":  accessToken,
		"refreshToken": refreshToken,
		"user":         newUserResponseData,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
