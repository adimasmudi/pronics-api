package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bidang struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	KategoriId primitive.ObjectID `json:"kategori_id" bson:"kategori_id"`
	LayananId []primitive.ObjectID `json:"layanan_id" bson:"layanan_id"`
	NamaBidang string `json:"nama_bidang"`
	CreatedById primitive.ObjectID `json:"created_by_id,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}