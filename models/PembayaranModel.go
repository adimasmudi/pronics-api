package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Pembayaran struct {
	ID                    primitive.ObjectID   `json:"id" bson:"_id"`
	OrderId            primitive.ObjectID   `json:"order_id" bson:"order_id"`
	Jenis string `json:"jenis"`
	PembayarId            primitive.ObjectID   `json:"pembayar_id" bson:"pembayar_id"`
	TerbayarId            primitive.ObjectID   `json:"terbayar_id" bson:"terbayar_id"`
	Status           string               `json:"status"`
	CreatedAt             time.Time            `json:"created_at"`
	UpdatedAt             time.Time            `json:"updated_at"`
}