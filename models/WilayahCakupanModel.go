package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WilayahCakupan struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	NamaWilayah string `json:"nama_wilayah"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}