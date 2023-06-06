package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Customer struct {
	UserId primitive.ObjectID `json:"user_id" bson:"user_id"`
	Username string `json:"username,omitempty"`
	GambarCustomer string `json:"gambarCustomer,omitempty"`
	AlamatCustomer []string `json:"alamat_customer,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}