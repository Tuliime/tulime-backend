package adverts

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

var PostAdvert = func(c *fiber.Ctx) error {
	advert := models.Advert{}
	store := models.Store{}
	user := c.Locals("user")

	if err := c.BodyParser(&advert); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if advert.UserID == "" || advert.ProductName == "" || advert.ProductDescription == "" {
		return fiber.NewError(fiber.StatusBadRequest,
			"Missing StoreID/UserID/ProductName/ProductDescription!")
	}

	// Get store id, if it doesn't exist, create one
	if advert.StoreID == "" {
		savedStore, err := store.FindByUser(advert.UserID, 1, "")
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if len(savedStore) > 0 {
			advert.StoreID = savedStore[0].ID
		} else {
			currentUser, ok := user.(models.User)
			if !ok {
				log.Printf("Invalid user type received: %+v", currentUser)
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			newStore, err := store.Create(models.Store{UserID: advert.UserID, Name: currentUser.Name})
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			advert.StoreID = newStore.ID
		}
	}

	newAdvert, err := advert.Create(advert)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert created successfully!",
		"data":    newAdvert,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
