package inputs

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddLayananInput struct {
	BidangId primitive.ObjectID `json:"bidang_id" binding:"required"`
	NamaLayanan   string `json:"nama_layanan" binding:"required"`
	Harga float64 `json:"harga" binding:"required"`
	AvailableTakeDelivery bool `json:"available_take_delivery" binding:"required"`
}