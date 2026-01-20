package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/Rcoleto-coder/go-api/internal/database"
	"github.com/Rcoleto-coder/go-api/internal/handlers"
	"github.com/Rcoleto-coder/go-api/internal/middleware"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	database.Connect(os.Getenv("MONGO_URI"))

	mux := http.NewServeMux()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)

	protected := middleware.AuthMiddleware(jwtSecret)
	mux.Handle("/", protected(http.HandlerFunc(handlers.Home)))

	handler := middleware.CORS(mux)

	log.Println("HTTP API running on http://localhost:" + port)
	http.ListenAndServe(":"+port, handler)

}
