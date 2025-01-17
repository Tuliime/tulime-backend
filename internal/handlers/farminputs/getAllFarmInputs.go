package farminputs

import (
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetAllFarmInputs = func(c *fiber.Ctx) error {
	farmInputs := models.FarmInputs{}
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

	allFarmInputs, err := farmInputs.FindAll(limit, categoryParam, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var prevCursor string
	if len(allFarmInputs) > 0 {
		prevCursor = allFarmInputs[len(allFarmInputs)-1].ID
	}

	pagination := map[string]interface{}{
		"limit":      limit,
		"prevCursor": prevCursor,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       allFarmInputs,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
