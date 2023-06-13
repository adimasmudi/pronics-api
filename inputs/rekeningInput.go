package inputs

type UpdateRekeningInput struct {
	NamaPemilikRekening string `json:"nama_pemilik_rekening" binding:"required"`
	IdBank              string `json:"id_bank" binding:"required"`
	NomerRekening       string `json:"nomer_rekening" binding:"required"`
}