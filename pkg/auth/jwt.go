package auth

import (
	"errors"
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

func GenerateRefreshToken(userID string, secret string) (string, error) {
	claims := jwt.MapClaims{
		// sub means subject, we pass the user ID in the token info
		"sub": userID, // Add user ID to the refresh token
		// the expiration time for the refresh token - 2 days
		"exp": time.Now().Add(2 * 24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateRefreshToken validates the refresh token and returns the user ID if valid.
func ValidateRefreshToken(tokenString, secret string) (string, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userID, ok := claims["sub"].(string)
	if !ok || userID == "" {
		return "", errors.New("user id not found in token")
	}

	return userID, nil
}
