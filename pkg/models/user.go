package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type User struct {
	ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Email         string             `bson:"email" json:"email"`
	Password      string             `bson:"password" json:"-"`
	Role          string             `bson:"role" json:"role"`
	CreatedAt     time.Time          `bson:"created_at" json:"created_at"`
	RefreshTokens []RefreshToken     `bson:"refresh_tokens" json:"-"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type RefreshToken struct {
	Token     string    `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}