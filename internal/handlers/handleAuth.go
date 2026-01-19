package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Rcoleto-coder/go-api/internal/database"
    "github.com/Rcoleto-coder/go-api/pkg/auth"
    "github.com/Rcoleto-coder/go-api/pkg/models"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var input models.User
	_ = json.NewDecoder(r.Body).Decode(&input)

	hash, _ := auth.HashPassword(input.Password)

	input.Password = hash
	input.CreatedAt = time.Now()
	input.Role = "user"

	collection := database.Client.
		Database(os.Getenv("DB_NAME")).
		Collection("users")

	_, err := collection.InsertOne(context.Background(), input)
	if err != nil {
		http.Error(w, "User exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var input models.User
	_ = json.NewDecoder(r.Body).Decode(&input)

	collection := database.Client.
		Database(os.Getenv("DB_NAME")).
		Collection("users")

	var user models.User
	err := collection.FindOne(
		context.Background(),
		map[string]string{"email": input.Email},
	).Decode(&user)

	if err != nil || auth.CheckPassword(user.Password, input.Password) != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, _ := auth.GenerateAccessToken(
		user.ID.Hex(),
		os.Getenv("JWT_SECRET"),
	)

	refreshToken, _ := auth.GenerateRefreshToken(
		os.Getenv("JWT_SECRET"),
	)

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
	})

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": accessToken,
	})
}
