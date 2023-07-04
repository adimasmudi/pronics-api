package formatters

import (
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AdminFormatter struct {
	ID          primitive.ObjectID `json:"id"`
	Username  string `json:"username"`
	Email       string `json:"email"`
}

type DashboardSummaryAdmin struct{
	TotalCustomer int `json:"total_customer"`
	TotalMitra int `json:"total_mitra"`
	TotalTransaksi int `json:"total_transaksi"`
}

func FormatAdmin(admin models.Admin) AdminFormatter {
	formatter := AdminFormatter{
		ID : admin.ID,
		Username:   admin.Username,
		Email:       admin.Email,
	}

	return formatter
}
