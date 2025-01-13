package vetdoctor

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetVetDoctorByUser = func(c *fiber.Ctx) error {
	vetDoctor := models.VetDoctor{}
	userID := c.Params("userId")

	vetDoctor, err := vetDoctor.FindByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   vetDoctor,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
