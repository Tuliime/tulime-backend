package middlewares

import (
	"log"
	"strings"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
)

func NonAuthCheck(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var bearerToken string
	user := models.User{}

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) > 1 {
			bearerToken = headerParts[1]
		}
	}

	if bearerToken == "" {
		user, err := user.FindByTelNumber(constants.AnonymousTelNumber)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}

		if user.ID == "" {
			log.Println("Anonymous User Doesn't Exist!")
		}
		c.Locals("userID", user.ID)
		c.Locals("user", user)
	}

	return c.Next()
}
