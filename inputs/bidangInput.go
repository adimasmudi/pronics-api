package inputs

import "go.mongodb.org/mongo-driver/bson/primitive"

type AddBidangInput struct {
	NamaBidang string `json:"nama_bidang"`
	KategoriId   primitive.ObjectID `json:"kategori_id"`
}