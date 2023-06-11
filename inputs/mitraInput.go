package inputs

type RegisterMitraInput struct {
	NamaLengkap         string `form:"nama_lengkap" binding:"required"`
	Email               string `form:"email" binding:"required"`
	Password            string `form:"password" binding:"required"`
	NamaToko            string `form:"nama_toko"`
	NomerTelepon        string `form:"no_telepon" binding:"required"`
	Deskripsi           string `form:"deskripsi"`
	Type                string `form:"type" binding:"required"`
	Alamat              string `form:"alamat" binding:"required"`
	MitraType           string `form:"type_mitra" binding:"required"`
	Wilayah             string `form:"wilayah" binding:"required"`
	Bidang              string `form:"bidang" binding:"required"`
	NamaPemilikRekening string `form:"nama_pemilik_rekening" binding:"required"`
	NamaBank            string `form:"nama_bank" binding:"required"`
	NomerRekening       string `form:"nomer_rekening" binding:"required"`
}

type UpdateProfilMitraInput struct{
	NamaLengkap  string `form:"nama_lengkap" binding:"required"`
	Email        string `form:"email" binding:"required"`
	Username     string `form:"username" binding:"required"`
	NoHandphone  string `form:"no_handphone"`
	Deskripsi    string `form:"bio"`
	JenisKelamin string `form:"jenis_kelamin"`
	TanggalLahir string `form:"tanggal_lahir"`
}
