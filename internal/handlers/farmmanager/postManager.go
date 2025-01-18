package farmmanager

import (
	"fmt"
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type PostFarmManagerInput struct {
	Name      string `validate:"string"`
	Gender    string `validate:"string"`
	RegNo     string `validate:"string"`
	Email     string `validate:"string"`
	TelNumber int    `validate:"number"`
}

var PostFarmManager = func(c *fiber.Ctx) error {
	farmManager := models.FarmManager{}
	userID := c.Params("userID")
	fmt.Println("userID:", userID)

	if err := c.BodyParser(&farmManager); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	farmManager.UserID = userID

	var input PostFarmManagerInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	savedFarmManager, err := farmManager.FindByUser(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmManager.ID != "" {
		return fiber.NewError(fiber.StatusBadRequest, "User already registered as a farm manager!")
	}

	farmManagerID, err := savedFarmManager.Create(farmManager)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	newFarmManager, err := farmManager.FindOne(farmManagerID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Your Farm Manager Profile has been created!",
		"data":    newFarmManager,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
