package adverts

import (
	"fmt"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

// TODO: To add image dimensions and compression
var PostAdvertImage = func(c *fiber.Ctx) error {
	advertID := c.Params("id")
	advert := models.Advert{}
	advertImage := models.AdvertImage{AdvertID: advertID}

	savedAdvert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	multipartForm, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString("Error parsing form")
	}

	fileHeaders := multipartForm.File["files"]

	var advertImages []models.AdvertImage

	// Validate image sizes (10 MB limit)
	const maxFileSize = 10 << 20 // 10 MB in bytes
	for _, fileHeader := range fileHeaders {
		if fileHeader.Size > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest,
				fmt.Sprintf("%s exceeds the 10 MB limit", fileHeader.Filename))
		}

	}

	// upload images
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

		advertImages = append(advertImages,
			models.AdvertImage{AdvertID: savedAdvert.ID,
				URL: imageUrl, Path: filePath})

	}

	if len(advertImages) > 0 {
		advertImages, err = advertImage.CreateMany(advertImages)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert Images uploaded successfully!",
		"data":    advertImages,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
