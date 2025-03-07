package device

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var EnableDevice = func(c *fiber.Ctx) error {
	device := models.Device{}
	deviceID := c.Params("id")

	if deviceID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Please provide deviceID!")
	}

	savedDevice, err := device.FindOne(deviceID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedDevice.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Device of provided id does'nt exist!")
	}

	if !savedDevice.NotificationDisabled {
		return fiber.NewError(fiber.StatusBadRequest, "Device is already enabled!")
	}

	savedDevice.NotificationDisabled = false

	updateDevice, err := savedDevice.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Device enabled successfully!",
		"data":    updateDevice,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
