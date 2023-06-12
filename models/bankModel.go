package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Bank struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	NamaBank string `json:"nama_bank"`
	LogoBank string `json:"logo_bank"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}