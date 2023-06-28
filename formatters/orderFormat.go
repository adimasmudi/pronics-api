package formatters

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderResponse struct {
	ID primitive.ObjectID `json:"id"`
	UserId primitive.ObjectID `json:"user_id"`
	MitraId primitive.ObjectID `json:"mitra_id"`
	TanggalOrder time.Time `json:"tanggal_order"`
	Status string `json:"status"`
	OrderDetail OrderDetailResponse `json:"order_detail"`
	TerakhirDiUpdate time.Time `json:"last_update"`
}