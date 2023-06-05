package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rekening struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	BankId primitive.ObjectID `json:"bank_id" bson:"bank_id"`
	NamaPemilik string `json:"nama_pemilik"`
	NomerRekening string `json:"nomer_rekening"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}