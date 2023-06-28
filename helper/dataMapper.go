package helper

import (
	"pronics-api/formatters"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapperCustomer(user models.User, customer models.Customer, alamat []formatters.AlamatResponse) (formatters.CustomerResponse){
	var data formatters.CustomerResponse
	var userData formatters.UserResponse

	userData.ID = user.ID
	userData.NamaLengkap = user.NamaLengkap
	userData.Email = user.Email
	userData.NoTelepon = user.NoTelepon
	userData.Bio = user.Deskripsi
	userData.JenisKelamin = user.JenisKelamin
	userData.TanggalLahir = user.TanggalLahir

	data.ID = customer.ID
	data.Username = customer.Username
	data.Alamat = alamat
	data.User = userData
	data.GambarUser = customer.GambarCustomer

	return data
}

func MapperMitra(user models.User, mitra models.Mitra, wilayah models.WilayahCakupan, bidang []models.Bidang)(formatters.MitraResponse){
	var data formatters.MitraResponse
	var userData formatters.UserResponse

	userData.ID = user.ID
	userData.NamaLengkap = user.NamaLengkap
	userData.Email = user.Email
	userData.NoTelepon = user.NoTelepon
	userData.Bio = user.Deskripsi
	userData.JenisKelamin = user.JenisKelamin
	userData.TanggalLahir = user.TanggalLahir

	data.ID = mitra.ID
	data.MitraType = mitra.MitraType
	data.NamaToko = mitra.NamaToko
	data.Alamat = mitra.Alamat
	data.Status = mitra.Status
	data.User = userData
	data.Wilayah = wilayah
	data.Bidang = bidang

	return data
}

func MapperOrder(customerId primitive.ObjectID, mitraId primitive.ObjectID, order models.Order, orderDetail formatters.OrderDetailResponse)(formatters.OrderResponse){
	var orderData formatters.OrderResponse
	orderData.ID = order.ID
	orderData.CustomerId = customerId
	orderData.MitraId = mitraId
	orderData.Status = order.Status
	orderData.OrderDetail = orderDetail
	orderData.TanggalOrder = order.TanggalOrderSelesai
	orderData.TerakhirDiUpdate = order.UpdatedAt

	return orderData
}