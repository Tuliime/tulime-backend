package adverts

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostAdvertImage = func(c *fiber.Ctx) error {
	advertID := c.Params("id")
	advert := models.Advert{}
	isPrimaryImageStr := c.FormValue("isPrimary")
	advertImage := models.AdvertImage{AdvertID: advertID}

	savedAdvert, err := advert.FindOne(advertID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if savedAdvert.ID == "" {
		return fiber.NewError(fiber.StatusNotFound, "Advert of provided id is not found!")
	}

	if isPrimaryImageStr == "" {
		isPrimaryImageStr = "false"
	}

	isPrimary, err := strconv.ParseBool(isPrimaryImageStr)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
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

	filePath = packages.GenFilePath(fileHeader.Filename)
	firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

	imageUrl, err = firebaseStorage.Add(file, fileHeader)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	advertImage.URL = imageUrl
	advertImage.Path = filePath
	advertImage.IsPrimary = isPrimary

	newAdvertImage, err := advertImage.Update()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert Image uploaded successfully!",
		"data":    newAdvertImage,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
