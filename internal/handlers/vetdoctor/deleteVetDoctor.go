package vetdoctor

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var DeleteVetDoctor = func(c *fiber.Ctx) error {
	vetDoctor := models.VetDoctor{}
	vetDoctorID := c.Params("id")

	savedVetDoctor, err := vetDoctor.FindOne(vetDoctorID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedVetDoctor.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Vet Doctor of provided id is not found!")
	}

	if err := vetDoctor.Delete(vetDoctorID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Vet Doctor deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
