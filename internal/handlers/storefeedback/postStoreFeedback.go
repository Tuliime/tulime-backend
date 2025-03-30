package storefeedback

import (
	"fmt"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostStoreFeedback = func(c *fiber.Ctx) error {
	feedback := models.StoreFeedback{}
	feedbackFile := models.StoreFeedbackFile{}

	storeID := c.Params("id")
	userID, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid userID type!")
	}
	feedback.UserID = userID
	feedback.StoreID = storeID

	feedback.Experience = c.FormValue("experience")
	feedback.Title = c.FormValue("title")
	feedback.Description = c.FormValue("description")

	var fileUploaded bool = true

	if feedback.Title == "" || feedback.Description == "" || feedback.Experience == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing title/description/experience!")
	}

	multipartForm, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form")
	}

	fileHeaders := multipartForm.File["files"]

	var feedbackFiles []models.StoreFeedbackFile

	// Validate image sizes (10 MB limit)
	if fileUploaded {
		const maxFileSize = 10 << 20 // 10 MB in bytes
		for _, fileHeader := range fileHeaders {
			if fileHeader.Size > maxFileSize {
				return fiber.NewError(fiber.StatusBadRequest,
					fmt.Sprintf("%s exceeds the 10 MB limit", fileHeader.Filename))
			}

		}
	}

	newFeedback, err := feedback.Create(feedback)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// upload images
	if fileUploaded {
		for _, fileHeader := range fileHeaders {
			file, err := fileHeader.Open()
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}
			defer file.Close()

			filePath := packages.GenFilePath(fileHeader.Filename)
			firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

			imageUrl, err := firebaseStorage.Add(file, fileHeader)
			if err != nil {
				return fiber.NewError(fiber.StatusBadRequest, err.Error())
			}

			feedbackFiles = append(feedbackFiles,
				models.StoreFeedbackFile{StoreFeedbackID: newFeedback.ID,
					URL: imageUrl, Path: filePath})

		}
	}

	if len(feedbackFiles) > 0 {
		createdFeedbackFiles, err := feedbackFile.CreateMany(feedbackFiles)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newFeedback.File = createdFeedbackFiles
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Feedback submitted successfully!",
		"data":    newFeedback,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
