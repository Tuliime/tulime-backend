package farminputs

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostFarmInputs = func(c *fiber.Ctx) error {
	farmInputs := models.FarmInputs{}
	farmInputs.Name = c.FormValue("name")
	farmInputs.Category = c.FormValue("category")
	farmInputs.Purpose = c.FormValue("purpose")
	farmInputs.Price, _ = strconv.ParseFloat(c.FormValue("price"), 64)
	farmInputs.PriceCurrency = c.FormValue("priceCurrency")
	farmInputs.Source = c.FormValue("source")
	farmInputs.SourceUrl = c.FormValue("sourceUrl")

	if farmInputs.Name == "" || farmInputs.Category == "" || farmInputs.Purpose == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing name/category/purpose!")
	}

	if farmInputs.Price == 0 || farmInputs.PriceCurrency == "" || farmInputs.Source == "" || farmInputs.SourceUrl == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing price/priceCurrency/source/sourceUrl!")
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

	farmInputs.ImageUrl = imageUrl
	farmInputs.ImagePath = filePath

	farmInput, err := farmInputs.Create(farmInputs)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Farm input created successfully!",
		"data":    farmInput,
	}

	return c.Status(fiber.StatusCreated).JSON(response)

}
