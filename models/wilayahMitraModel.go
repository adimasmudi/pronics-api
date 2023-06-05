package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type WilayahMitra struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	WilayahId primitive.ObjectID `json:"wilayah_id" bson:"wilayah_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}