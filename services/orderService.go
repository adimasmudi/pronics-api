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
	GetAllOrder(ctx context.Context)([]formatters.OrderResponse, error)
}

type orderService struct{
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
	customerRepository repositories.CustomerRepository
	orderRepository repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	orderPaymentRepository repositories.OrderPaymentRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewOrderService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository,orderPaymentRepository repositories.OrderPaymentRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *orderService{
	return &orderService{userRepository, mitraRepository, customerRepository, orderRepository,orderDetailRepository,orderPaymentRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
}

func (s *orderService) CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error){
	var orderData formatters.OrderResponse
	var orderDetailData formatters.OrderDetailResponse

	var orderPaymentToDisplay models.OrderPayment
	var layananData formatters.LayananDetailMitraResponse
	var bidangData formatters.BidangResponse

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

		if err == nil{

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


			layananToOrder, err := s.layananRepository.GetById(ctx, orderDetailToDisplay.LayananId)

			if err != nil{
				layananMitraToOrder, err := s.layananMitraRepository.GetById(ctx, orderDetailToDisplay.LayananId)

				if err != nil{
					return orderData, err
				}

				layananData.ID = layananMitraToOrder.ID
				layananData.NamaLayanan = layananMitraToOrder.NamaLayanan
				layananData.Harga = layananMitraToOrder.Harga
			}else{
				layananData.ID = layananToOrder.ID
				layananData.NamaLayanan = layananToOrder.NamaLayanan
				layananData.Harga = layananToOrder.Harga
			}

			orderPaymentToDisplay, err = s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetailToDisplay.ID)

			if err != nil{
				return orderData, err
			}
		}

		orderPaymentData := helper.MapperOrderPayment(orderPaymentToDisplay)

		orderDetailData = helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPaymentData)
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

func (s *orderService) GetAllOrder(ctx context.Context)([]formatters.OrderResponse, error){
	var orderResponses []formatters.OrderResponse

	var bidangData formatters.BidangResponse
	var layananData formatters.LayananDetailMitraResponse

	orders, err := s.orderRepository.GetAllOrder(ctx)

	if err != nil{
		return orderResponses, err
	}

	for _, order := range orders{
		orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

		if err != nil{
			return orderResponses, err
		}

		bidang, err := s.bidangRepository.GetById(ctx, orderDetail.BidangId)

		if err != nil{
			return orderResponses, err
		}

		kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

		if err != nil{
			return orderResponses, err
		}

		bidangData.ID = bidang.ID
		bidangData.Kategori = kategori.NamaKategori
		bidangData.NamaBidang = bidang.NamaBidang


		layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)

		if err != nil{
			layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

			if err != nil{
				return orderResponses, err
			}

			layananData.ID = layananMitra.ID
			layananData.NamaLayanan = layananMitra.NamaLayanan
			layananData.Harga = layananMitra.Harga
		}else{
			layananData.ID = layanan.ID
			layananData.NamaLayanan = layanan.NamaLayanan
			layananData.Harga = layanan.Harga
		}

		orderPayment, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetail.ID)

		if err != nil{
			return orderResponses, err
		}

		orderPaymentMapping := helper.MapperOrderPayment(orderPayment)
		orderDetailMapping := helper.MapperOrderDetail(orderDetail,bidangData,layananData,orderPaymentMapping)
		orderMapping := helper.MapperOrder(order.CustomerId, order.MitraId,order,orderDetailMapping)

		orderResponses = append(orderResponses, orderMapping)
	}

	return orderResponses, nil


}