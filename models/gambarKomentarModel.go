package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GambarKomentar struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	KomentarId primitive.ObjectID `json:"komentar_id" bson:"komentar_id"`
	Gambar string `json:"gambar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}