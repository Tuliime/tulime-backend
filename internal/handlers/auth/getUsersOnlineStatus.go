package auth

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
)

var GetUsersOnlineStatus = func(c *fiber.Ctx) error {
	userIDListEncoding := c.Query("userIDListEncoding")
	log.Printf("userIDListEncoding %+v :", userIDListEncoding)

	decodedBytes, err := base64.StdEncoding.DecodeString(userIDListEncoding)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	userIDListStr := string(decodedBytes)
	fmt.Println("Decoded JSON:", userIDListStr)
	if userIDListStr == "" {
		return fiber.NewError(fiber.StatusInternalServerError, "Provided encoding without content!")
	}

	var userIDList []string
	if userIDListStr != "" {
		err = json.Unmarshal(decodedBytes, &userIDList)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid format! encoded data must array of strings.")
		}
	}
	fmt.Printf("userIDList: %v\n", userIDList)

	// fetch users online status here

	response := fiber.Map{
		"status": "success",
		"data":   userIDList,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
