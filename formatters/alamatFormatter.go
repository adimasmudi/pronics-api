package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type AlamatResponse struct {
	ID     primitive.ObjectID `json:"id"`
	Alamat string             `json:"alamat"`
}