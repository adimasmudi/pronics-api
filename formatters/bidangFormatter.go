package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type BidangResponse struct {
	ID         primitive.ObjectID
	NamaBidang string `json:"nama_bidang"`
}