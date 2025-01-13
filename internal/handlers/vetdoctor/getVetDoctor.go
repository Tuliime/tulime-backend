package vetdoctor

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var GetVetDoctor = func(c *fiber.Ctx) error {
	vetDoctor := models.VetDoctor{}
	vetDoctorID := c.Params("id")

	vetDoctor, err := vetDoctor.FindOne(vetDoctorID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status": "success",
		"data":   vetDoctor,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
