package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type BidangResponse struct {
	ID         primitive.ObjectID
	Kategori string `json:"kategori"`
	NamaBidang string `json:"nama_bidang"`
}

type DetailBidangResponse struct{
	ID         primitive.ObjectID
	NamaBidang string `json:"nama_bidang"`
	Kategori string `json:"kategori"`
	Layanan []LayananResponse `json:"layanan"`
}