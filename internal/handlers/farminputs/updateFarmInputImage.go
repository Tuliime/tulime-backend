package farminputs

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateFarmInputImage = func(c *fiber.Ctx) error {
	farmInputID := c.Params("id")
	farmInputs := models.FarmInputs{ID: farmInputID}

	savedFarmInput, err := farmInputs.FindOne(farmInputID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedFarmInput.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Farm of provided id is not found!")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

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
	savedFilePath := savedFarmInput.ImagePath

	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err := firebaseStorage.Update(fileReader, file, savedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFarmInput.ImageUrl = imageUrl
	savedFarmInput.ImagePath = filePath

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
