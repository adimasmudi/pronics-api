package inputs

type UpdateProfilCustomerInput struct {
	NamaLengkap  string    `form:"nama_lengkap" binding:"required"`
	Email        string    `form:"email" binding:"required"`
	Username     string    `form:"username" binding:"required"`
	NoHandphone  string    `form:"no_handphone"`
	Deskripsi    string    `form:"bio"`
	JenisKelamin string    `form:"jenis_kelamin"`
	TanggalLahir string `form:"tanggal_lahir"`
}