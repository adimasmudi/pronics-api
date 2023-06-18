package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type LayananMitra struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	BidangId primitive.ObjectID `json:"bidang_id" bson:"bidang_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	NamaLayanan string `json:"nama_layanan"`
	Harga float64 `json:"harga"`
	AvailableTakeDelivery bool `json:"available_take_delivery"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}