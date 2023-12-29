package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID                primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Email             string             `json:"email" bson:"email" validate:"required,email"`
	Password          string             `json:"password" bson:"password" binding:"required,min=8"`
	FirstName         string             `json:"firstName" bson:"firstName" validate:"required"`
	LastName          string             `json:"lastName" bson:"lastName" validate:"required"`
	Phone             string             `json:"phone" bson:"phone" validate:"required"`
	Role              string             `json:"role" bson:"role" validate:"required,oneof=ADMIN USER"`
	VerificationCode  string             `json:"verificationCode" bson:"verificationCode"`
	ResetPasswordCode string             `json:"resetPasswordCode" bson:"resetPasswordCode"`
	// Token             string             `json:"token,omitempty" bson:"token,omitempty"`
	// RefreshToken      string             `json:"refreshToken,omitempty" bson:"refreshToken,omitempty"`
	CreatedAt time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt time.Time `json:"updated_at" bson:"updated_at"`
}

type SignInInput struct {
	Email    string `json:"email" bson:"email" binding:"required"`
	Password string `json:"password" bson:"password" binding:"required"`
}
