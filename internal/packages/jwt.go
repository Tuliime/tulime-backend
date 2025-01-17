package packages

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

func SignJWTToken(userID string, tokenType string) (string, error) {
	var tokenExpiration int64

	if tokenType == "accessToken" {
		tokenExpiration = time.Now().Add(12 * time.Hour).Unix()
	} else if tokenType == "refreshToken" {
		tokenExpiration = time.Now().Add(24 * 30 * time.Hour).Unix()
	} else {
		return "", fmt.Errorf("invalid token type: %s", tokenType)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": userID,
		"exp":    tokenExpiration,
		"iat":    time.Now().Unix(),
	})

	secretKey := os.Getenv("JWT_SECRET")
	var jwtSecretKey = []byte(secretKey)

	accessToken, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign JWT token: %v", err)
	}

	return accessToken, nil
}
