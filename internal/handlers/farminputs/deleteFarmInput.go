package farminputs

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var DeleteFarmInput = func(c *fiber.Ctx) error {
	farmInputs := models.FarmInputs{}
	farmInputID := c.Params("id")

	savedFarmInput, err := farmInputs.FindOne(farmInputID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmInput.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Product of provided id is not found!")
	}

	if err := farmInputs.Delete(farmInputID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFilePath := savedFarmInput.ImagePath
	firebaseStorage := packages.FirebaseStorage{}

	if err := firebaseStorage.Delete(savedFilePath); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Farm input deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
