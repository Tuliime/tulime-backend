package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var DeleteAgroProduct = func(c *fiber.Ctx) error {
	agroProduct := models.Agroproduct{}

	agroProductID := c.Params("id")

	savedAgroProduct, err := agroProduct.FindOne(agroProductID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedAgroProduct.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Product of provided id is not found!")
	}

	if err := agroProduct.Delete(agroProductID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var savedFilePath string

	// Extract savedFilePath from the url if saved imagePath is null
	if packages.HasImagePath(savedAgroProduct) {
		savedFilePath = savedAgroProduct.ImagePath
	} else {
		savedFilePath, err = packages.ExtractFilePath(savedAgroProduct.ImageUrl)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	firebaseStorage := packages.FirebaseStorage{}

	if err := firebaseStorage.Delete(savedFilePath); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Product deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
