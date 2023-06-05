package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `json:"id" bson:"_id"`
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon"`
	Password string `json:"password"`
	Deskripsi string `json:"deskripsi"`
	JenisKelamin string `json:"jenis_kelamin"`
	Type string `json:"type"` // customer / mitra
	TanggalLahir time.Time `json:"tanggal_lahir"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}