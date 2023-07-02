package formatters

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderResponse struct {
	ID primitive.ObjectID `json:"id"`
	TransaksiId string `json:"transaksi_id"`
	Customer CustomerResponse `json:"customer"`
	Mitra MitraResponse `json:"mitra"`
	TanggalOrder time.Time `json:"tanggal_order"`
	Status string `json:"status"`
	OrderDetail OrderDetailResponse `json:"order_detail"`
	TerakhirDiUpdate time.Time `json:"last_update"`
}

type OrderHistoryResponse struct{
	ID primitive.ObjectID `json:"id"`
	Name string `json:"nama"`
	Layanan string `json:"layanan"`
	Status string `json:"status"`
	AlamatPemesanan string `json:"alamat_pemesanan"`
	TotalBayar int `json:"total_bayar"`
}