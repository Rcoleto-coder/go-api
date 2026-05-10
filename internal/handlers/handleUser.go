package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/Rcoleto-coder/go-api/internal/middleware"
	"github.com/Rcoleto-coder/go-api/internal/services"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" || id == r.URL.Path {
		http.Error(w, "user id required", http.StatusBadRequest)
		return
	}

	authUserID, ok := middleware.UserIDFromContext(r.Context())
	if !ok || authUserID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if id != authUserID {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	user, err := services.GetUserByID(r.Context(), id)
	if err != nil {
		switch {
		case err == services.ErrInvalidUserID:
			http.Error(w, "invalid user id", http.StatusBadRequest)
		case err == services.ErrUserNotFound:
			http.Error(w, "user not found", http.StatusNotFound)
		default:
			log.Println("GetUserByID error:", err)
			http.Error(w, "internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		log.Println("GetUser encode error:", err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
