package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MitraResponse struct {
	ID primitive.ObjectID `json:"id"`
	User UserResponse `json:"user_data"`
	MitraType string `json:"mitra_type"`
	NamaToko string `json:"nama_toko"`
	Alamat string `json:"alamat"`
	Status string `json:"status"`
	Wilayah models.WilayahCakupan `json:"wilayah"`
	Bidang []models.Bidang `json:"bidang"`
}

type KatalogResponse struct{
	ID primitive.ObjectID `json:"id"`
	Gambar string `json:"gambar"`
	Name string `json:"name"`
	RangePrice string `json:"range_price"`
	Bidang []BidangResponse `json:"bidang"`
	Rating int `json:"rating"`
}