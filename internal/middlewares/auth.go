package middlewares

import (
	"fmt"
	"os"
	"strings"

	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
)

func Auth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	var bearerToken string

	if authHeader != "" && strings.HasPrefix(authHeader, "Bearer") {
		headerParts := strings.SplitN(authHeader, " ", 2)
		if len(headerParts) > 1 {
			bearerToken = headerParts[1]
		}
	}

	if bearerToken == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "You are not logged in! Please log in to gain access.")
	}

	secretKey := os.Getenv("JWT_SECRET")
	var jwtSecretKey = []byte(secretKey)

	token, err := jwt.Parse(bearerToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return fiber.NewError(fiber.StatusForbidden, err.Error())
	}

	var userID string
	claims, validJWTClaim := token.Claims.(jwt.MapClaims)
	if !validJWTClaim || !token.Valid {
		return fiber.NewError(fiber.StatusForbidden, "Invalid token. Please log in again.")
	}

	if userIDClaim, ok := claims["userId"].(string); ok {
		userID = userIDClaim
	}

	User := models.User{}
	user, err := User.FindOne(userID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if user.ID == "" {
		return fiber.NewError(fiber.StatusForbidden, "The user belonging to this token no longer exists!")
	}

	c.Locals("userID", userID)

	return c.Next()
}
