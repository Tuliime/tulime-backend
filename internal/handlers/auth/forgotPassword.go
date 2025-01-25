package auth

import (
	"fmt"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var ForgotPassword = func(c *fiber.Ctx) error {
	user := models.User{}

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	user, err := user.FindByTelNumber(user.TelNumber)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "We couldn't find user with provided telephone number!")
	}

	otp := models.OTP{UserID: user.ID}
	optCode, err := otp.Create(otp)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	fmt.Println("OTP:", optCode)
	//TODO:To Send an sms message here

	response := map[string]interface{}{
		"status":  "success",
		"message": "OTP has been sent to your number",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
