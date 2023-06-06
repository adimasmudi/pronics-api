package inputs

type RegisterMitraInput struct {
	NamaLengkap         string `json:"nama_lengkap" binding:"required"`
	Email               string `json:"email" binding:"required"`
	Password            string `json:"password" binding:"required"`
	NamaToko            string `json:"nama_toko"`
	NomerTelepon        string `json:"no_telepon" binding:"required"`
	Deskripsi           string `json:"deskripsi"`
	Alamat              string `json:"alamat" binding:"required"`
	MitraType           string `json:"type_mitra" binding:"required"`
	Wilayah             string `json:"wilayah" binding:"required"`
	Bidang              string `json:"bidang" binding:"required"`
	NamaPemilikRekening string `json:"nama_pemilik_rekening" binding:"required"`
	NamaBank            string `json:"nama_bank" binding:"required"`
	NomerRekening       string `json:"nomer_rekening" binding:"required"`
}
