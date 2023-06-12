package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserFormatter struct {
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoHandphone string `json:"no_handphone"`
	Tipe string `json:"tipe"`
}

func FormatUser(user models.User) UserFormatter {
	formatter := UserFormatter{
		NamaLengkap : user.NamaLengkap,
		Email : user.Email ,
		NoHandphone: user.NoTelepon,
		Tipe : user.Type,
	}

	return formatter
}

type UserResponse struct{
	ID primitive.ObjectID `json:"id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Bio string `json:"bio"`
	JenisKelamin string `json:"jenis_kelamin"`
	TanggalLahir string `json:"tanggal_lahir"`
}