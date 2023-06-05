package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderDetail struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	OrderId primitive.ObjectID `json:"order_id" bson:"order_id"`
	BidangId primitive.ObjectID `json:"bidang_id" bson:"bidang_id"`
	JenisOrder string `json:"jenis_order"`
	Merk string `json:"merk"`
	JenisLayanan string `json:"jenis_layanan"`
	DeskripsiKerusakan string `json:"deskripsi_kerusakan"`
	AlamatPemesanan string `json:"alamat_pemesanan"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}