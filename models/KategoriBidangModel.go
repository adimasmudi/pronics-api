package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Kategori struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	BidangId []primitive.ObjectID `json:"bidang_id,omitempty" bson:"bidang_id"`
	NamaKategori string `json:"nama_bidang"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}