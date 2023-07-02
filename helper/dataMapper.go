package helper

import (
	"pronics-api/formatters"
	"pronics-api/models"
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

func MapperOrder(customer formatters.CustomerResponse, mitra formatters.MitraResponse, order models.Order, orderDetail formatters.OrderDetailResponse)(formatters.OrderResponse){
	var orderData formatters.OrderResponse
	orderData.ID = order.ID
	orderData.TransaksiId = order.TransaksiId
	orderData.Customer = customer
	orderData.Mitra = mitra
	orderData.Status = order.Status
	orderData.OrderDetail = orderDetail
	orderData.TanggalOrder = order.TanggalOrderSelesai
	orderData.TerakhirDiUpdate = order.UpdatedAt

	return orderData
}

func MapperOrderDetail(orderDetail models.OrderDetail, bidang formatters.BidangResponse, layanan formatters.LayananDetailMitraResponse, orderPayment formatters.OrderPaymentResponse) formatters.OrderDetailResponse{
	var orderData formatters.OrderDetailResponse

	orderData.ID = orderDetail.ID
	orderData.Bidang = bidang
	orderData.JenisOrder = orderDetail.JenisOrder
	orderData.AlamatPesanan = orderDetail.AlamatPemesanan
	orderData.DeskripsiKerusakan = orderDetail.DeskripsiKerusakan
	orderData.JenisOrder = orderDetail.JenisOrder
	orderData.LastUpdate = orderDetail.UpdatedAt
	orderData.Layanan = layanan
	orderData.Merk = orderDetail.Merk
	orderData.OrderPayment = orderPayment

	return orderData
}

func MapperOrderPayment(orderPayment models.OrderPayment, jarak float64) formatters.OrderPaymentResponse{
	var orderData formatters.OrderPaymentResponse

	orderData.ID = orderPayment.ID
	orderData.BiayaAplikasi = orderPayment.BiayaAplikasi
	orderData.BiayaPelayanan = orderPayment.BiayaPelayanan
	orderData.BiayaPerjalanan = orderPayment.BiayaPerjalanan
	orderData.BuktiBayar = orderPayment.BuktiBayar
	orderData.Diskon = orderPayment.Diskon
	orderData.TotalBiaya = orderPayment.TotalBiaya
	orderData.Jarak = jarak
	orderData.MetodePembayaran = orderPayment.MetodePembayaran
	orderData.LastUpdate = orderPayment.UpdatedAt

	return orderData

}