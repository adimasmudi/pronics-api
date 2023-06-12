package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AlamatCustomer struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Alamat string `json:"alamat"`
	IsUtama bool `json:"is_utama"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}