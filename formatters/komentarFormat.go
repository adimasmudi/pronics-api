package formatters

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type KomentarResponse struct {
	ID primitive.ObjectID `json:"id"`
	FotoCustomer string `json:"foto_customer"`
	NamaCustomer string `json:"nama_customer"`
	RatingGiven float64 `json:"rating_given"`
	Tanggal time.Time `json:"tanggal"`
	Layanan string `json:"layanan_dipesan"`
	Gambar []string `json:"gambar_komentar"`
	TotalSuka int `json:"suka"`
	Komentar string `json:"komentar"`

}

type KomentarDetailMitraResponse struct{
	OverallRating float64 `json:"overall_rating"`
	RatingCount int `json:"rating_count"`
	CommentCount int `json:"comment_count"`
	AllKomentar []KomentarResponse `json:"comments"`
}