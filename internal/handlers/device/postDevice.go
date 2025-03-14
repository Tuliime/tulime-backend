package device

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// TODO: allow same device token for multiple users
var PostDevice = func(c *fiber.Ctx) error {
	device := models.Device{}

	if err := c.BodyParser(&device); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if device.Name == "" || device.UserID == "" || device.Token == "" || device.TokenType == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/userId/token/tokenType!")
	}

	savedDevice, err := device.FindByTokenAndUser(device.Token, device.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedDevice.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "This device is already added!")
	}

	newDevice, err := device.Create(device)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Device created successfully!",
		"data":    newDevice,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
