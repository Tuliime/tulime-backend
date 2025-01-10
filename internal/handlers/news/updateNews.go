package news

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

type UpdateNewsInput struct {
	Title    string `validate:"string"`
	Category string `validate:"string"`
	Source   string `validate:"string"`
}

var UpdateNews = func(c *fiber.Ctx) error {
	news := models.News{}

	if err := c.BodyParser(&news); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var input UpdateNewsInput
	errors := packages.ValidateInput(c, &input)
	if len(errors) > 0 {
		log.Printf("Validation Error %+v :", errors)
		// TODO: Implement channels to send error detail to the default
		// fiber error handler
		return fiber.NewError(fiber.StatusBadRequest, "Validation Error")
	}

	newsID := c.Params("id")

	savedNews, err := news.FindOne(newsID)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedNews.Title = news.Title
	savedNews.Category = news.Category
	savedNews.Description = news.Description

	updatedNews, err := savedNews.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Updated successfully!",
		"data":    updatedNews,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
