package auth

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateUserInput struct {
	Name      string `validate:"string"`
	TelNumber string `validate:"telephoneNumber"`
}

var UpdateUser = func(c *fiber.Ctx) error {
	user := models.User{}
	userID := c.Params("id")

	if err := c.BodyParser(&user); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateUserInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	savedUser, err := user.FindOne(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedUser.Name = user.Name
	savedUser.TelNumber = user.TelNumber

	updatedUser, err := savedUser.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedUser,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
