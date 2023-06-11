package services

import (
	"context"
	"fmt"
	"os"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerService interface {
	GetCustomerProfile(ctx context.Context, ID primitive.ObjectID) (formatters.CustomerResponse, error)
	UpdateProfileCustomer(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilCustomerInput, fileName string) (*mongo.UpdateResult, error)
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

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, user.ID)

	if err != nil {
		return data, err
	}

	data = helper.MapperCustomer(user, customer)

	return data, nil
}

func (s *customerService) UpdateProfileCustomer(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilCustomerInput, fileName string) (*mongo.UpdateResult, error){
	newCustomer := bson.M{
		"username" : input.Username,
		"gambarcustomer": os.Getenv("CLOUD_STORAGE_READ_LINK")+"customer/"+fileName,
		"updatedat" : time.Now(),
	}

	newUser := bson.M{
		"namalengkap" : input.NamaLengkap,
		"email" : input.Email,
		"notelepon" : input.NoHandphone,
		"deskripsi" : input.Deskripsi,
		"jeniskelamin" : input.JenisKelamin,
		"tanggallahir" : input.TanggalLahir,
		"updatedat": time.Now(),
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx,ID)

	if err != nil{
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, ID, newUser)

	if err != nil{
		return nil, err
	}

	updatedCustomer, err := s.customerRepository.UpdateProfil(ctx, customer.ID,newCustomer)

	if err != nil{
		return nil, err
	}

	fmt.Println(updatedCustomer)

	return updatedUser, nil
}