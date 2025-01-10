package news

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var DeleteNews = func(c *fiber.Ctx) error {
	news := models.News{}
	newsID := c.Params("id")

	savedNews, err := news.FindOne(newsID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if savedNews.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Product of provided id is not found!")
	}

	if err := news.Delete(newsID); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFilePath := savedNews.ImagePath
	firebaseStorage := packages.FirebaseStorage{}

	if err := firebaseStorage.Delete(savedFilePath); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "News deleted successfully!",
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
