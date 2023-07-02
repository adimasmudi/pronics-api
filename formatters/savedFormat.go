package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type SavedResponse struct {
	ID primitive.ObjectID `json:"id"`
	Mitra KatalogResponse `json:"mitra"`
}