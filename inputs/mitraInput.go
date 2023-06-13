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
	IdBank              string `form:"id_bank" binding:"required"`
	NomerRekening       string `form:"nomer_rekening" binding:"required"`
}

type UpdateProfilMitraInput struct {
	NamaLengkap  string `form:"nama_lengkap" binding:"required"`
	Email        string `form:"email" binding:"required"`
	NoHandphone  string `form:"no_handphone" binding:"required"`
	Deskripsi    string `form:"deskripsi"`
	JenisKelamin string `form:"jenis_kelamin"`
	TanggalLahir string `form:"tanggal_lahir"`
	NamaToko     string `form:"nama_toko"`
	Alamat       string `form:"alamat" binding:"required"`
}
