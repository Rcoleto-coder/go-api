package services

import (
	"context"
	"errors"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Rcoleto-coder/go-api/internal/database"
	"github.com/Rcoleto-coder/go-api/pkg/models"
)

var (
	ErrInvalidUserID = errors.New("invalid user id")
	ErrUserNotFound  = errors.New("user not found")
)

func GetUserByID(ctx context.Context, id string) (*models.UserResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, ErrInvalidUserID
	}

	collection := database.Client.
		Database(os.Getenv("DB_NAME")).
		Collection("users")

	var user models.User
	err = collection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &models.UserResponse{
		ID:    user.ID.Hex(),
		Email: user.Email,
		Role:  user.Role,
	}, nil
}
