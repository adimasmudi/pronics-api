package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Admin struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Username string `json:"username"`
	Email string `json:"email"`
	Password string `json:"password"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}