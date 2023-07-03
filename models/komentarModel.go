package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Komentar struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	OrderId primitive.ObjectID `json:"order_id" bson:"order_id"`
	GambarKomentar []string `json:"gambar_komentar" bson:"gambar_komentar"`
	Penyuka []primitive.ObjectID `json:"penyuka"`
	Rating float64 `json:"rating"`
	Komentar string `json:"komentar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}