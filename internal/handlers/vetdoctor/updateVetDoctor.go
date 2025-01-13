package vetdoctor

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateVetDoctorInput struct {
	Name          string `validate:"string"`
	Gender        string `validate:"string"`
	LicenseNumber string `validate:"string"`
	Email         string `validate:"string"`
	TelNumber     int    `validate:"number"`
}

var UpdateVetDoctor = func(c *fiber.Ctx) error {
	vetDoctor := models.VetDoctor{}
	vetDoctorID := c.Params("id")

	if err := c.BodyParser(&vetDoctor); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateVetDoctorInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	savedVetDoctor, err := vetDoctor.FindOne(vetDoctorID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedVetDoctor.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Vet Doctor of provided id is not found!")
	}

	// TODO: To validate emails, gender
	savedVetDoctor.Name = vetDoctor.Name
	savedVetDoctor.Email = vetDoctor.Email
	savedVetDoctor.Gender = vetDoctor.Gender
	savedVetDoctor.LicenseNumber = vetDoctor.LicenseNumber
	savedVetDoctor.TelNumber = vetDoctor.TelNumber

	updatedVetDoctor, err := savedVetDoctor.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Your Vet Doctor Profile has been updated!",
		"data":    updatedVetDoctor,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
