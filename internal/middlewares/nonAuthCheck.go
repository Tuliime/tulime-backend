package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
)

func NonAuthCheck(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var bearerToken string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) > 1 {
			bearerToken = headerParts[1]
		}
	}

	if bearerToken == "" {
		c.Locals("userID", "")
	}

	return c.Next()
}
