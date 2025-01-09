package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostAgroProduct = func(c *fiber.Ctx) error {
	agroProduct := models.Agroproduct{}

	agroProduct.Name = c.FormValue("name")
	agroProduct.Category = c.FormValue("category")

	if agroProduct.Name == "" || agroProduct.Category == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/category!")
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
	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err := firebaseStorage.Add(fileReader, file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	agroProduct.ImageUrl = imageUrl
	agroProduct.ImagePath = filePath

	newAgroProduct, err := agroProduct.Create(agroProduct)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Agro product created successfully!",
		"data":    newAgroProduct,
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}
