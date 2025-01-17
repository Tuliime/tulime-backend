package packages

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// SignJWTToken fn signs a new jwt token with userID and returns it
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

// DecodeJWT fn decodes the provided jwt token and returns the userID
func DecodeJWT(JWTToken string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")
	var jwtSecretKey = []byte(secretKey)

	token, err := jwt.Parse(JWTToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return "", err
	}

	var userID string
	claims, validJWTClaim := token.Claims.(jwt.MapClaims)
	if !validJWTClaim || !token.Valid {
		return "", errors.New("invalid Token")
	}

	if userIDClaim, ok := claims["userId"].(string); ok {
		userID = userIDClaim
	}

	return userID, nil
}
