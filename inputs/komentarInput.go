package inputs

type KomentarInput struct {
	Rating   float64 `form:"rating" json:"rating" binding:"required"`
	Komentar string  `form:"komentar" json:"komentar" binding:"required"`
}