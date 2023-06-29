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
	orderDetailRepository repositories.OrderDetailRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewOrderService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *orderService{
	return &orderService{userRepository, mitraRepository, customerRepository, orderRepository,orderDetailRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
}

func (s *orderService) CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error){
	var orderData formatters.OrderResponse
	var orderDetailData formatters.OrderDetailResponse

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return orderData, err
	}

	mitra, err := s.mitraRepository.GetMitraById(ctx, mitraId)

	if err != nil{
		return orderData, err
	}

	order, err := s.orderRepository.GetOrderTemporaryByCustomerIdNMitraId(ctx, customer.ID, mitra.ID)


	if err == nil{
		orderDetailToDisplay, err := s.orderDetailRepository.GetByOrderId(ctx,order.ID)

		if err != nil{
			return orderData, err
		}

		var bidangData formatters.BidangResponse

		bidangToOrder, err := s.bidangRepository.GetById(ctx, orderDetailToDisplay.BidangId)

		if err != nil{
			return orderData, err
		}

		kategoriToOrder, err := s.kategoriRepository.GetById(ctx, bidangToOrder.KategoriId)

		if err != nil{
			return orderData, err
		}

		bidangData.ID = bidangToOrder.ID
		bidangData.Kategori = kategoriToOrder.NamaKategori
		bidangData.NamaBidang = bidangToOrder.NamaBidang

		var layananData formatters.LayananResponse

		layananToOrder, err := s.layananRepository.GetById(ctx, orderDetailToDisplay.LayananId)

		if err != nil{
			layananMitraToOrder, err := s.layananMitraRepository.GetById(ctx, orderDetailToDisplay.LayananId)

			if err != nil{
				return orderData, err
			}

			layananData.ID = layananMitraToOrder.ID
			layananData.NamaLayanan = layananMitraToOrder.NamaLayanan
		}else{
			layananData.ID = layananToOrder.ID
			layananData.NamaLayanan = layananToOrder.NamaLayanan
		}

		var orderPayment formatters.OrderPaymentResponse // sementara

		orderDetailData = helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPayment)
		orderData = helper.MapperOrder(customer.ID, mitra.ID, order, orderDetailData)

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


	orderData = helper.MapperOrder(customer.ID, mitra.ID, orderAdded, orderDetailData)

	return orderData, nil
}
