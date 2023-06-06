package models

import (
	"time"
)

type User struct {
	NamaLengkap string `json:"nama_lengkap"`
	Email string `json:"email"`
	NoTelepon string `json:"no_telepon,omitempty"`
	Password string `json:"password"`
	Deskripsi string `json:"deskripsi,omitempty"`
	JenisKelamin string `json:"jenis_kelamin,omitempty"`
	Type string `json:"type"` // customer / mitra
	TanggalLahir time.Time `json:"tanggal_lahir,omitempty"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}