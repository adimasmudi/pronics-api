package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderPayment struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	OrderDetailId primitive.ObjectID `json:"order_detail_id" bson:"order_detail_id"`
	BiayaPelayanan int64 `json:"biaya_pelayanan"`
	BiayaPerjalanan int64 `json:"biaya_perjalanan"`
	Diskon int64 `json:"diskon"`
	BiayaAplikasi int64 `json:"biaya_aplikasi"`
	TotalBiaya int64 `json:"total_biaya"`
	BiayaAfterDiskon int64 `json:"biaya_after_diskon"`
	MetodePembayaran string `json:"metode_pembayaran"`
	BuktiBayar string `json:"bukti_bayar"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}