package services

import (
	"context"
	"pronics-api/formatters"
	"pronics-api/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CustomerService interface {
	GetCustomerProfile(ctx context.Context, ID primitive.ObjectID) (formatters.CustomerResponse, error)
}

type customerService struct {
	userRepository     repositories.UserRepository
	customerRepository repositories.CustomerRepository
}

func NewCustomerService(userRepository repositories.UserRepository, customerRepository repositories.CustomerRepository) *customerService{
	return &customerService{userRepository, customerRepository}
}

func (s *customerService) GetCustomerProfile(ctx context.Context, ID primitive.ObjectID) (formatters.CustomerResponse, error){ 
	var data formatters.CustomerResponse
	var userData formatters.UserResponse

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, user.ID)

	if err != nil {
		return data, err
	}

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

	return data, nil
}