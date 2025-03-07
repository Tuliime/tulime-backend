package device

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetDeviceByUser = func(c *fiber.Ctx) error {
	device := models.Device{}
	userID := c.Params("userID")
	if userID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide userID")
	}

	devices, err := device.FindByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   devices,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
