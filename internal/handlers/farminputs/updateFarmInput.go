package farminputs

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

// type FarmInputs struct {
// 	ID        string    `gorm:"column:id;type:uuid;primaryKey" json:"id"`
// 	Name      string    `gorm:"column:name;unique;not null;index" json:"name"`
// 	Purpose   string    `gorm:"column:purpose;not null" json:"purpose"`
// 	Category  string    `gorm:"column:category;not null;index" json:"category"`
// 	ImageUrl  string    `gorm:"column:imageUrl;not null" json:"imageUrl"`
// 	ImagePath string    `gorm:"column:imagePath;not null" json:"imagePath"`
// 	CreatedAt time.Time `gorm:"column:createdAt;index" json:"createdAt"`
// 	UpdatedAt time.Time `gorm:"column:updatedAt;index" json:"updatedAt"`
// }

type UpdateFarmInputValidator struct {
	Name     string `validate:"string"`
	Category string `validate:"string"`
	Purpose  string `validate:"string"`
}

var UpdateFarmInput = func(c *fiber.Ctx) error {
	farmInput := models.FarmInputs{}

	if err := c.BodyParser(&farmInput); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateFarmInputValidator
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	farmInputID := c.Params("id")

	savedFarmInput, err := farmInput.FindOne(farmInputID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFarmInput.Name = farmInput.Name
	savedFarmInput.Category = farmInput.Category
	savedFarmInput.Purpose = farmInput.Purpose

	updatedFarmInput, err := savedFarmInput.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedFarmInput,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
