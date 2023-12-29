package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserSession struct {
	ID         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	UserID     primitive.ObjectID `json:"userId, bson:"userId" validate:"required"`
	Token      string             `json:"token" bson:"token" validate:"required"`
	DeviceID   string             `json:"devideId" bson:"devideId" validate:"required"`
	DeviceType string             `json:"deviceType" bson:"deviceType" validate:"required"`
}
