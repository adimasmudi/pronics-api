package services

import (
	"context"
	"errors"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderService interface {
	CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error)
}

type orderService struct{
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
	customerRepository repositories.CustomerRepository
	orderRepository repositories.OrderRepository
}

func NewOrderService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository) *orderService{
	return &orderService{userRepository, mitraRepository, customerRepository, orderRepository}
}

func (s *orderService) CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error){
	var orderData formatters.OrderResponse

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return orderData, err
	}

	mitra, err := s.mitraRepository.GetMitraById(ctx, mitraId)

	if err != nil{
		return orderData, err
	}

	order, err := s.orderRepository.GetOrderTemporaryByCustomerIdNMitraId(ctx, customer.ID, mitra.ID)

	var orderDetail formatters.OrderDetailResponse // sementara

	if err == nil{
		orderData = helper.MapperOrder(customer.ID, mitra.ID, order, orderDetail)

		return orderData, errors.New(constants.TemporaryOrderExistMessage)
	}

	newOrder := models.Order{
		ID : primitive.NewObjectID(),
		CustomerId: customer.ID,
		MitraId: mitra.ID,
		Status : constants.OrderTemporary,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	addOrder, err := s.orderRepository.Create(ctx, newOrder)

	if err != nil{
		return orderData, err
	}

	orderAdded, err := s.orderRepository.GetById(ctx,addOrder.InsertedID.(primitive.ObjectID))

	if err != nil{
		return orderData, err
	}

	orderData = helper.MapperOrder(customer.ID, mitra.ID, orderAdded, orderDetail)

	return orderData, nil
}
