package inputs

type RegisterUserInput struct {
	NamaLengkap string `json:"nama_lengkap" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Password    string `json:"password" binding:"required"`
	Type        string `json:"type" binding:"required"`
}

type LoginUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}