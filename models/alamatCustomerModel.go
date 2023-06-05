package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlamatCustomer struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Alamat string `json:"alamat"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}