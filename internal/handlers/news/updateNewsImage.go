package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var UpdateNewsImage = func(c *fiber.Ctx) error {
	newsID := c.Params("id")
	news := models.News{ID: newsID}

	savedNews, err := news.FindOne(newsID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedNews.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "News of provided id is not found!")
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
	savedFilePath := savedNews.ImagePath

	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err := firebaseStorage.Update(fileReader, file, savedFilePath)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedNews.ImageUrl = imageUrl
	savedNews.ImagePath = filePath

	updatedNews, err := savedNews.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Update created successfully!",
		"data":    updatedNews,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
