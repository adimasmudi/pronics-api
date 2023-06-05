package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ContactUs struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	NamaLengkap string `json:"nama_lengkap"`
	NomerHandphone string `json:"nomer_handphone"`
	Email string `json:"email"`
	Pesan string `json:"pesan"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}