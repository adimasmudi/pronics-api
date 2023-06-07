package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Mitra struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	GambarMitra []primitive.ObjectID `json:"gambar_mitra,omitempty"`
	Wilayah primitive.ObjectID `json:"wilayah"`
	Bidang []primitive.ObjectID `json:"bidang"`
	NamaToko string `json:"nama_toko,omitempty"`
	Alamat string `json:"alamat"`
	MitraType string `json:"mitra_type"` // individu or toko
	Status string `json:"status"` // active or inactive
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}