package formatters

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderDetailResponse struct {
	ID primitive.ObjectID `json:"id"`
	Bidang BidangResponse `json:"bidang"`
	JenisOrder string `json:"jenis_order"`
	Layanan LayananDetailMitraResponse `json:"layanan"`
	Merk string `json:"merk"`
	DeskripsiKerusakan string `json:"deskripsi_kerusakan"`
	AlamatPesanan string `json:"alamat_pesanan"`
	OrderPayment OrderPaymentResponse `json:"order_payment"`
	LastUpdate time.Time `json:"last_update"`
}