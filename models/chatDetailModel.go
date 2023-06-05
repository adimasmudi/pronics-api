package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ChatDetail struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	Chat_id primitive.ObjectID `json:"chat_id" bson:"chat_id"`
	Sender_id primitive.ObjectID `json:"sender_id" bson:"sender_id"`
	Receiver_id primitive.ObjectID `json:"receiver_id" bson:"receiver_id"`
	Pesan string `json:"pesan"`
	Media string `json:"media"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}