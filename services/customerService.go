package services

import (
	"context"
	"pronics-api/formatters"
	"pronics-api/helper"
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