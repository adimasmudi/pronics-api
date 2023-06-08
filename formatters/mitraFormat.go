package formatters

import "go.mongodb.org/mongo-driver/bson/primitive"

type MitraResponse struct {
	ID primitive.ObjectID `json:"id"`
	User UserResponse `json:"user_data"`
	MitraType string `json:"mitra_type"`
	NamaToko string `json:"nama_toko"`
	Alamat string `json:"alamat"`
	Status string `json:"status"`
}