package storefeedback

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var UpdateStoreFeedback = func(c *fiber.Ctx) error {
	feedback := models.StoreFeedback{}
	feedbackID := c.Params("feedbackID")

	if err := c.BodyParser(&feedback); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if feedback.Title == "" || feedback.Description == "" || feedback.Experience == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing title/description/experience!")
	}

	savedFeedback, err := feedback.FindOne(feedbackID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())

	}

	savedFeedback.Experience = feedback.Experience
	savedFeedback.Title = feedback.Title
	savedFeedback.Description = feedback.Description

	updatedFeedback, err := savedFeedback.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())

	}

	response := fiber.Map{
		"status":  "success",
		"message": "Feedback updated successfully!",
		"data":    updatedFeedback,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
