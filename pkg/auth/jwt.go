package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func GenerateAccessToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		// sub means subject, we pass the user ID in the token info 
		"sub": userID, 
		// the expiration time for the access token - 15 minutes
		"exp": time.Now().Add(15 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func GenerateRefreshToken(secret string) (string, error) {
	claims := jwt.MapClaims{
		// the expiration time for the refresh token - 2 days
		"exp": time.Now().Add(2 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
