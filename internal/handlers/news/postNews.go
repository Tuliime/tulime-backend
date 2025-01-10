package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostNews = func(c *fiber.Ctx) error {
	news := models.News{}

	news.Title = c.FormValue("title")
	news.Category = c.FormValue("category")
	news.Source = c.FormValue("source")
	// TODO: consider manipulating the date string "postedAt"
	// news.PostedAt = c.FormValue("postedAt")

	if news.Title == "" || news.Category == "" || news.Source == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing title/category/source!")
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
	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err := firebaseStorage.Add(fileReader, file)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	news.ImageUrl = imageUrl
	news.ImagePath = filePath

	createdNews, err := news.Create(news)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "News created successfully!",
		"data":    createdNews,
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}
