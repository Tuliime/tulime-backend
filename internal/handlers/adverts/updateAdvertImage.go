package adverts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateAdvertImage = func(c *fiber.Ctx) error {
	advertImage := models.AdvertImage{}
	advertImageID := c.Params("advertImageID")

	savedAdvertImage, err := advertImage.FindOne(advertImageID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvertImage.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert Image of provided id is not found!")
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var filePath string
	var imageUrl string

	// Validate file size (10 MB limit)
	const maxFileSize = 10 << 20 // 10 MB in bytes
	if fileHeader.Size > maxFileSize {
		return fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 10 MB limit")
	}
	file, err := fileHeader.Open()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	defer file.Close()

	filePath = packages.GenFilePath(fileHeader.Filename)
	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err = firebaseStorage.Update(file, fileHeader, savedAdvertImage.Path)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedAdvertImage.URL = imageUrl
	savedAdvertImage.Path = filePath

	updatedAdvertImage, err := savedAdvertImage.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert Image updated successfully!",
		"data":    updatedAdvertImage,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
