package formatters

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderPaymentResponse struct {
	ID              primitive.ObjectID `json:"id"`
	BiayaPelayanan  float64            `json:"biaya_pelayanan"`
	BiayaPerjalanan float64            `json:"biaya_perjalanan"`
	Diskon float64 `json:"diskon"`
	BiayaAplikasi float64 `json:"biaya_aplikasi"`
	TotalBiaya float64 `json:"total_biaya"`
	MetodePembayaran string `json:"metode_pembayaran"`
	BuktiBayar string `json:"bukti_bayar"`
	LastUpdate time.Time `json:"last_update"`
}