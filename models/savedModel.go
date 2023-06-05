package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Saved struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Customer_id primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	Mitra_id primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}