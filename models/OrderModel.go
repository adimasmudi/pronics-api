package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	KomentarId []primitive.ObjectID `json:"komentar_id" bson:"komentar_id"`
	TanggalOrder time.Time `json:"tanggal_order"`
	Status string `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}