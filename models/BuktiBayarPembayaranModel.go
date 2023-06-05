package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BuktiBayarPembayaran struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	PembayaranId primitive.ObjectID `json:"pembayaran_id" bson:"pembayaran_id"`
	BuktiBayar string `json:"bukti_bayar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}