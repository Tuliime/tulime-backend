package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var SignUp = func(c *fiber.Ctx) error {
	user := models.User{}
	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.Name == "" || user.TelNumber == 0 || user.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing username/telNumber/password!")
	}

	savedUser, err := user.FindByTelNumber(user.TelNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedUser.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "Telephone number already registered!")
	}

	err = user.SetRole("user")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userId, err := user.Create(user)

	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	accessToken, err := packages.SignJWTToken(userId)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	newUser := map[string]interface{}{
		"id":        userId,
		"name":      user.Name,
		"telNumber": user.TelNumber,
		"role":      user.Role,
	}
	response := map[string]interface{}{
		"status":      "success",
		"message":     "Account created successfully",
		"accessToken": accessToken,
		"user":        newUser,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
