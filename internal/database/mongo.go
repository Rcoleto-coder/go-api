package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB *mongo.Database
var Client *mongo.Client

func Connect(uri string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))

	// Check for connection error
	if err != nil {
		log.Fatal("Mongo connect error: ", err)
	}
	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal("Mongo ping failed: ", err)
	}
	// Use the database name "go-api"
	DB = client.Database("go-api")

	// Set the global client variable
	Client = client
	log.Println("MongoDB connected")
}
