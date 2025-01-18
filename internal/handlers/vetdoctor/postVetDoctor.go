package vetdoctor

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type PostVetDoctorInput struct {
	Name          string `validate:"string"`
	Gender        string `validate:"string"`
	LicenseNumber string `validate:"string"`
	Email         string `validate:"string"`
	TelNumber     int    `validate:"number"`
}

var PostVetDoctorManager = func(c *fiber.Ctx) error {
	vetDoctor := models.VetDoctor{}
	userID := c.Params("userID")

	if err := c.BodyParser(&vetDoctor); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	vetDoctor.UserID = userID

	var input PostVetDoctorInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	savedFarmManager, err := vetDoctor.FindByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmManager.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "User already registered as a vet!")
	}

	farmManagerID, err := savedFarmManager.Create(vetDoctor)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newFarmManager, err := vetDoctor.FindOne(farmManagerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Your Veterinary Doctor Profile has been created!",
		"data":    newFarmManager,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
