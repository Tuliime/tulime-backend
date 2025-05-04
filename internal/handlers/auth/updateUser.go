package auth

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateUserInput struct {
	Name        string `validate:"string"`
	TelNumber   int    `validate:"number"` //TODO: To validate telNumber in "ValidateInput" in  a better way
	Email       string `validate:"string"`
	Gender      string `validate:"string"`
	DateOfBirth string `validate:"string"`
	Country     string `validate:"string"`
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
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// TODO:To validate tel phone on every update
	savedUser.Name = user.Name
	savedUser.TelNumber = user.TelNumber

	if savedUser.Email != user.Email {
		emailUser, err := user.FindByEmail(user.Email)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if emailUser.ID != "" {
			return fiber.NewError(fiber.StatusBadRequest, "Can't update to already existing email!")
		}
		savedUser.Email = user.Email
	}

	if user.Gender != "" {
		savedUser.Gender = user.Gender
	}
	if user.DateOfBirth != "" {
		savedUser.DateOfBirth = user.DateOfBirth
	}
	if user.Country != "" {
		savedUser.Country = user.Country
	}

	updatedUser, err := savedUser.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	updateUserMap := fiber.Map{
		"id":             updatedUser.ID,
		"name":           updatedUser.Name,
		"telNumber":      updatedUser.TelNumber,
		"role":           updatedUser.Role,
		"imageUrl":       updatedUser.ImageUrl,
		"profileBgColor": updatedUser.ProfileBgColor,
		"chatroomColor":  updatedUser.ChatroomColor,
		"createdAt":      updatedUser.CreatedAt,
		"updatedAt":      updatedUser.UpdatedAt,
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updateUserMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
