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

	database.Connect(os.Getenv("MONGO_URI"))

	mux := http.NewServeMux()
	mux.HandleFunc("/register", handlers.Register)
	mux.HandleFunc("/login", handlers.Login)

	protected := middleware.AuthMiddleware(os.Getenv("JWT_SECRET"))
	mux.Handle("/home", protected(http.HandlerFunc(handlers.Home)))

	log.Println("HTTPS API running on https://localhost:8000")

	http.ListenAndServeTLS(
		":8000",
		os.Getenv("TLS_CERT"),
		os.Getenv("TLS_KEY"),
		mux,
	)
}
