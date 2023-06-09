package formatters

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KategoriWithBidangResponse struct {
	ID         primitive.ObjectID `json:"id"`
	NamaKategori string `json:"nama_kategori"`
	Bidang       []BidangResponse     `json:"bidang"`
}
