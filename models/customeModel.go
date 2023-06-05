package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username string `json:"username"`
	AlamatCustomer []string `json:"alamat_customer"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}