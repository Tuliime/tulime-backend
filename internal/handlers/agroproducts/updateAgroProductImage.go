package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateAgroProductImage = func(c *fiber.Ctx) error {
	agroProductID := c.Params("id")
	agroProduct := models.Agroproduct{ID: agroProductID}

	savedAgroProduct, err := agroProduct.FindOne(agroProductID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedAgroProduct.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Product of provided id is not found!")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// Validate file size (10 MB limit)
	const maxFileSize = 10 << 20 // 10 MB in bytes
	if file.Size > maxFileSize {
		return fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 10 MB limit")
	}

	fileReader, err := file.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	defer fileReader.Close()

	filePath := packages.GenFilePath(file.Filename)
	var savedFilePath string

	// Extract savedFilePath from the url if saved imagePath is null
	if agroProduct.HasImagePath(savedAgroProduct) {
		savedFilePath = savedAgroProduct.ImagePath
	} else {
		savedFilePath, err = packages.ExtractFilePath(savedAgroProduct.ImageUrl)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err := firebaseStorage.Update(fileReader, file, savedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedAgroProduct.ImageUrl = imageUrl
	savedAgroProduct.ImagePath = filePath

	updatedAgroProduct, err := savedAgroProduct.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Agro product created successfully!",
		"data":    updatedAgroProduct,
	}

	return c.Status(fiber.StatusOK).JSON(response)

}
