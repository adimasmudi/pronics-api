package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RekeningResponse struct {
	ID primitive.ObjectID `json:"id"`
	Bank models.Bank `json:"bank"`
	NamaPemilik string `json:"nama_pemilik"`
	NomerRekening string `json:"nomer_rekening"`
}