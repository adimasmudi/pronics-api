package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type LayananResponse struct {
	ID         primitive.ObjectID
	NamaLayanan string `json:"nama_layanan"`
}