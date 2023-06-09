package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon,omitempty"`
	Password string `json:"password,omitempty"`
	Deskripsi string `json:"deskripsi,omitempty"`
	JenisKelamin string `json:"jenis_kelamin,omitempty"`
	Type string `json:"type"` // customer / mitra
	TanggalLahir time.Time `json:"tanggal_lahir,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}