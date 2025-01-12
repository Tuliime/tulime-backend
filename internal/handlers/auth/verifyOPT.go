package auth

import (
	"time"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var VerifyOTP = func(c *fiber.Ctx) error {
	otp := models.OTP{}

	if err := c.BodyParser(&otp); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	otpInput := otp.OTP

	savedOtp, err := otp.FindByOTP(otpInput)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedOtp.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid OTP!")
	}

	if savedOtp.IsUsed || savedOtp.IsVerified || savedOtp.ExpiresAt.Before(time.Now()) {
		return fiber.NewError(fiber.StatusBadRequest, "OPT is already used or invalid or expired!")
	}

	savedOtp.IsVerified = true
	otp, err = savedOtp.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	otpMap := map[string]interface{}{
		"isVerified": otp.IsVerified,
		"otp":        otpInput,
	}

	response := map[string]interface{}{
		"status":  "success",
		"message": "OTP verified successfully",
		"otp":     otpMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
