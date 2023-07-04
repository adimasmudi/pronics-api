package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type LayananResponse struct {
	ID         primitive.ObjectID
	NamaLayanan string `json:"nama_layanan"`
}

type LayananDetailMitraResponse struct{
	ID primitive.ObjectID `json:"id"`
	BidangId primitive.ObjectID `json:"bidang_id"`
	NamaLayanan string `json:"nama_layanan"`
	Harga float64 `json:"harga"`
}