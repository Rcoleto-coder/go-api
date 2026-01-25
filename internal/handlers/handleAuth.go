package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"io"
	"strings"

	"github.com/Rcoleto-coder/go-api/internal/database"
	"github.com/Rcoleto-coder/go-api/pkg/auth"
	"github.com/Rcoleto-coder/go-api/pkg/models"
	"github.com/Rcoleto-coder/go-api/pkg/utils"
)

type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest

	// Debug: log raw body
	rawBody, _ := io.ReadAll(r.Body)
	log.Println("Raw request body:", string(rawBody))
	r.Body = io.NopCloser(strings.NewReader(string(rawBody))) // Reset body for decoder

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		log.Println("Invalid JSON in Register request", err)
		return
	}

	// Normalize
	req.Email = utils.NormalizeEmail(req.Email)
	req.Password = utils.NormalizePassword(req.Password)

	// Validate
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		log.Println("Email and password are required")
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, "password must be at least 6 characters", http.StatusBadRequest)
		log.Println("password must be at least 6 characters")
		return
	}

	hash, err := auth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "failed to hash password", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Email:     req.Email,
		Password:  hash,
		CreatedAt: time.Now(),
		Role:      "user",
	}

	collection := database.Client.
		Database(os.Getenv("DB_NAME")).
		Collection("users")

	_, err = collection.InsertOne(context.Background(), user)
	if err != nil {
		http.Error(w, "user already exists", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	// Normalize
	req.Email = utils.NormalizeEmail(req.Email)
	req.Password = utils.NormalizePassword(req.Password)

	// Validate
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	collection := database.Client.
		Database(os.Getenv("DB_NAME")).
		Collection("users")

	var user models.User
	err := collection.FindOne(
		context.Background(),
		map[string]string{"email": req.Email},
	).Decode(&user)

	if err != nil || auth.CheckPassword(user.Password, req.Password) != nil {
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(
		user.ID.Hex(),
		os.Getenv("JWT_SECRET"),
	)
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	refreshToken, err := auth.GenerateRefreshToken(user.ID.Hex(), os.Getenv("JWT_SECRET"))
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "refreshToken",
		Value:    refreshToken,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/refresh",
	})

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": accessToken,
	})
}

func Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refreshToken")
	if err != nil {
		http.Error(w, "refresh token missing", http.StatusUnauthorized)
		return
	}

	userID, err := auth.ValidateRefreshToken(cookie.Value, os.Getenv("JWT_SECRET"))
	if err != nil {
		http.Error(w, "invalid refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, err := auth.GenerateAccessToken(userID, os.Getenv("JWT_SECRET"))
	if err != nil {
		http.Error(w, "token generation failed", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"accessToken": accessToken,
	})
}
