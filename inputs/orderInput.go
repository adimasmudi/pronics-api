package inputs

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddOrUpdateOrderDetailInput struct {
	BidangId primitive.ObjectID `json:"bidang_id" binding:"required"`
	Merk string `json:"merk" binding:"required"`
	LayananId primitive.ObjectID `json:"layanan_id" binding:"required"`
	DeskripsiKerusakan string `json:"deskripsi_kerusakan" binding:"required"`
	AlamatPesanan string `json:"alamat_pemesanan" binding:"required"`
}