package auth

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateUserImage = func(c *fiber.Ctx) error {
	user := models.User{}

	user, err := user.FindOne(c.Params("id"))
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}
	if user.ID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "We couldn't user find of the provided id")
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
	var imageUrl string

	if user.ImageUrl == "" {
		// Add Image to firebase storage when no imageUrl
		imageUrl, err = firebaseStorage.Add(fileReader, file)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	} else {
		// Update Image in firebase storage when has imageUrl
		imageUrl, err = firebaseStorage.Update(fileReader, file, user.ImagePath)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	user.ImageUrl = imageUrl
	user.ImagePath = filePath

	user, err = user.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userMap := map[string]interface{}{
		"id":        user.ID,
		"imageUrl":  user.ImageUrl,
		"updatedAt": user.UpdatedAt,
	}

	response := fiber.Map{
		"status":  "success",
		"message": "User image updated successfully!",
		"data":    userMap,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
