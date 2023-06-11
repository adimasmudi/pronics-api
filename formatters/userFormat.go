package formatters

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserResponse struct{
	ID primitive.ObjectID `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Bio string `json:"bio"`
	JenisKelamin string `json:"jenis_kelamin"`
	TanggalLahir string `json:"tanggal_lahir"`
}