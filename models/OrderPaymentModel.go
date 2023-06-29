package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderPayment struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	OrderDetailId primitive.ObjectID `json:"order_detail_id" bson:"order_detail_id"`
	BiayaPelayanan float64 `json:"biaya_pelayanan"`
	BiayaPerjalanan float64 `json:"biaya_perjalanan"`
	Diskon float64 `json:"diskon,omitempty"`
	BiayaAplikasi float64 `json:"biaya_aplikasi"`
	TotalBiaya float64 `json:"total_biaya"`
	BiayaAfterDiskon float64 `json:"biaya_after_diskon,omitempty"`
	MetodePembayaran string `json:"metode_pembayaran,omitempty"`
	BuktiBayar string `json:"bukti_bayar,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}