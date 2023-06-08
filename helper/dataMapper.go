package helper

import (
	"pronics-api/formatters"
	"pronics-api/models"
)

func MapperCustomer(user models.User, customer models.Customer) (formatters.CustomerResponse){
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
	data.Alamat = customer.AlamatCustomer // sementara, harusnya bisa get alamat customer
	data.User = userData
	data.GambarUser = customer.GambarCustomer

	return data
}

func MapperMitra(user models.User, mitra models.Mitra)(formatters.MitraResponse){
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

	return data
}