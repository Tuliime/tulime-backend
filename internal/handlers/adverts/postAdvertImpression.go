package adverts

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/Tuliime/tulime-backend/internal/handlers/ipinfo"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

// type AdvertIDs struct {
// 	AdvertIDList string `json:"advertIDList"`
// }

var PostAdvertImpression = func(c *fiber.Ctx) error {
	advertImpression := models.AdvertImpression{}
	location := models.Location{}
	advertIDs := AdvertIDs{}
	device := c.Get("X-Device")

	userID, ok := c.Locals("userID").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid userID type!")
	}

	clientIP, ok := c.Locals("clientIP").(string)
	if !ok {
		return fiber.NewError(fiber.StatusInternalServerError, "Invalid client type!")
	}

	if err := c.BodyParser(&advertIDs); err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if advertIDs.AdvertIDList == "" {
		return fiber.NewError(fiber.StatusBadRequest, "You provided an empty advertIDList!")
	}

	var advertIDList []string
	err := json.Unmarshal([]byte(advertIDs.AdvertIDList), &advertIDList)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest,
			"Invalid advertIDs format! Must be a JSON stringified array of strings.")
	}

	fmt.Printf("advertIDList:%+v\n", advertIDList)

	savedLocation, err := location.FindByIP(clientIP)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	location.ID = savedLocation.ID

	if savedLocation.ID == "" {
		// ip := "197.239.8.162"
		info, err := ipinfo.GetIPInfo(clientIP)
		if err != nil {
			log.Fatalf("Error fetching IP info: %v", err)
		}
		fmt.Printf("Found ip location:%+v\n", info)

		ipinfoJson, err := json.Marshal(&info)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newLocation, err := location.Create(models.Location{UserID: userID,
			Info: models.JSONB(ipinfoJson)}, clientIP)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		location.ID = newLocation.ID
	}

	var advertImpressions []models.AdvertImpression

	for _, advertID := range advertIDList {
		advertImpressions = append(advertImpressions,
			models.AdvertImpression{AdvertID: advertID, UserID: userID,
				LocationID: location.ID, Device: device})
	}

	newAdvertImpressions, err := advertImpression.Create(advertImpressions)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	response := fiber.Map{
		"status":  "success",
		"message": "Advert impressions created successfully!",
		"data":    newAdvertImpressions,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
