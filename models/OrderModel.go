package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	TransaksiId string `json:"transaksi_id,omitempty"`
	CustomerId primitive.ObjectID `json:"customer_id" bson:"customer_id"`
	MitraId primitive.ObjectID `json:"mitra_id" bson:"mitra_id"`
	KomentarId []primitive.ObjectID `json:"komentar_id,omitempty" bson:"komentar_id"`
	TanggalOrderSelesai time.Time `json:"tanggal_order,omitempty"`
	Status string `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}