package agroproducts

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

var GetAllProducts = func(c *fiber.Ctx) error {
	filePath := filepath.Join("internal", "data", "agroProducts.json")

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error reading products data: "+err.Error())
	}

	var response interface{}
	err = json.Unmarshal(data, &response)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing products data: "+err.Error())
	}

	return c.Status(fiber.StatusOK).JSON(response)
}