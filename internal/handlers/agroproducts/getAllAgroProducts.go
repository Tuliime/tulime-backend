package agroproducts

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

// var GetAllProducts = func(c *fiber.Ctx) error {
// 	filePath := filepath.Join("internal", "data", "agroProducts.json")

// 	data, err := os.ReadFile(filePath)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Error reading products data: "+err.Error())
// 	}

// 	var response interface{}
// 	err = json.Unmarshal(data, &response)
// 	if err != nil {
// 		return fiber.NewError(fiber.StatusInternalServerError, "Error parsing products data: "+err.Error())
// 	}

// 	return c.Status(fiber.StatusOK).JSON(response)
// }

var GetAllAgroProducts = func(c *fiber.Ctx) error {
	agroProduct := models.Agroproduct{}
	limitParam := c.Query("limit")
	categoryParam := c.Query("category")
	cursorParam := c.Query("cursor")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if cursorParam == "" {
		cursorParam = ""
	}
	if categoryParam == "" {
		categoryParam = ""
	}

	agroProducts, err := agroProduct.FindAll(limit, categoryParam, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	if len(agroProducts) > 0 {
		prevCursor = agroProducts[len(agroProducts)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       agroProducts,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
