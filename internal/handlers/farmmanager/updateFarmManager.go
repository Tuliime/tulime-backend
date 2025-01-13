package farmmanager

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateFarmManagerInput struct {
	Name      string `validate:"string"`
	Gender    string `validate:"string"`
	RegNo     string `validate:"string"`
	Email     string `validate:"string"`
	TelNumber int    `validate:"number"`
}

var UpdateFarmManager = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	farmManagerID := c.Params("id")

	if err := c.BodyParser(&farmManager); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateFarmManagerInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	savedFarmManager, err := farmManager.FindOne(farmManagerID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmManager.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Farm manager of provided id is not found!")
	}

	// TODO: To validate emails, gender
	savedFarmManager.Name = farmManager.Name
	savedFarmManager.Email = farmManager.Email
	savedFarmManager.Gender = farmManager.Gender
	savedFarmManager.RegNo = farmManager.RegNo
	savedFarmManager.TelNumber = farmManager.TelNumber

	updatedFarmManager, err := savedFarmManager.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Your Farm Manager Profile has been updated!",
		"data":    updatedFarmManager,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
