package store

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateStoreBgImage = func(c *fiber.Ctx) error {
	store := models.Store{}
	storeID := c.Params("id")

	savedStore, err := store.FindOne(storeID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedStore.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Store of provided id is not found!")
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

	imageProcessor := packages.ImageProcessor{}
	compressedFileBuf, err := imageProcessor.CompressMultipartFile(file, 75)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to compress image")
	}

	filePath = packages.GenFilePath(fileHeader.Filename)
	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	if savedStore.BackgroundImageUrl != "" {
		// Update existing image
		imageUrl, err = firebaseStorage.UpdateFromBuffer(compressedFileBuf,
			savedStore.BackgroundImagePath)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		// Add new image
		imageUrl, err = firebaseStorage.AddFromBuffer(compressedFileBuf)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	savedStore.BackgroundImageUrl = imageUrl
	savedStore.BackgroundImagePath = filePath

	updatedStore, err := savedStore.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Store updated successfully!",
		"data":    updatedStore,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
