package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MitraResponse struct {
	ID primitive.ObjectID `json:"id"`
	User UserResponse `json:"user_data"`
	MitraType string `json:"mitra_type"`
	NamaToko string `json:"nama_toko"`
	Alamat string `json:"alamat"`
	Status string `json:"status"`
	Wilayah models.WilayahCakupan `json:"wilayah"`
	Bidang []models.Bidang `json:"bidang"`
}

type KatalogResponse struct{
	ID primitive.ObjectID `json:"id"`
	Gambar string `json:"gambar"`
	Name string `json:"name"`
	MinPrice float64 `json:"minimal_price"`
	MaxPrice float64 `json:"maximal_price"`
	Distance float64 `json:"jarak"`
	Bidang []BidangResponse `json:"bidang"`
	Rating float64 `json:"rating"`
}
type DetailMitraResponse struct{
	ID primitive.ObjectID `json:"id"`
	GaleriImage []string `json:"galeri_image"`
	FotoProfil string `json:"foto_profil"`
	NamaToko string `json:"nama_toko"`
	NamaPemilik string `json:"nama_pemilik"`
	Deskripsi string `json:"deskripsi"`
	Bidang []string `json:"bidang"`
	Layanan []LayananDetailMitraResponse `json:"layanan"`
	Ulasan interface{} `json:"ulasan"`
}

type MitraDashboardSummaryResponse struct{
	ID primitive.ObjectID `json:"id"`
	NamaToko string `json:"nama_toko"`
	NamaPemilik string `json:"nama_pemilik"`
	Email string `json:"email"`
	NoHandphone string `json:"no_handphone"`
	JumlahTransaksiSelesai int `json:"jumlah_transaksi_selesai"`
	StatusMitra string `json:"status_mitra"`
}

type DashboardSummaryMitra struct{
	TotalOrderSelesai int `json:"total_order_selesai"`
	TotalPendapatanBersih float64 `json:"total_pendapatan_bersih"`
}