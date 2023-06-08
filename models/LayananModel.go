package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Layanan struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	BidangId primitive.ObjectID `json:"bidang_id" bson:"bidang_id"`
	NamaLayanan string `json:"nama_layanan"`
	Harga int64 `json:"harga"`
	AvailableTakeDelivery bool `json:"available_take_delivery"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}