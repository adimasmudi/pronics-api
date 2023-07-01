package services

import (
	"context"
	"errors"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderService interface {
	CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error)
	GetAllOrder(ctx context.Context)([]formatters.OrderResponse, error)
	GetOrderDetail(ctx context.Context, orderId primitive.ObjectID) (formatters.OrderResponse, error)
	UpdateStatusOrder(ctx context.Context, orderId primitive.ObjectID, input inputs.UpdateStatusOrderInput) (*mongo.UpdateResult, error)
	GetAllOrderMitra(ctx context.Context, userId primitive.ObjectID, status string) ([]formatters.OrderResponse, error)
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

	sort.Slice(orderResponses, func(i, j int) bool {
		return orderResponses[i].TerakhirDiUpdate.Unix() < orderResponses[j].TerakhirDiUpdate.Unix()
	})

	return orderResponses, nil
}

func (s *orderService) GetOrderDetail(ctx context.Context, orderId primitive.ObjectID) (formatters.OrderResponse, error){
	var orderResponse formatters.OrderResponse

	order, err := s.orderRepository.GetById(ctx, orderId)

	if err != nil{
		return orderResponse, err
	}

	var bidangData formatters.BidangResponse
	var layananData formatters.LayananDetailMitraResponse

	orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

		if err != nil{
			return orderResponse, err
		}

		bidang, err := s.bidangRepository.GetById(ctx, orderDetail.BidangId)

		if err != nil{
			return orderResponse, err
		}

		kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

		if err != nil{
			return orderResponse, err
		}

		bidangData.ID = bidang.ID
		bidangData.Kategori = kategori.NamaKategori
		bidangData.NamaBidang = bidang.NamaBidang


		layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)

		if err != nil{
			layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

			if err != nil{
				return orderResponse, err
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
			return orderResponse, err
		}

		orderPaymentMapping := helper.MapperOrderPayment(orderPayment)
		orderDetailMapping := helper.MapperOrderDetail(orderDetail,bidangData,layananData,orderPaymentMapping)
		orderResponse = helper.MapperOrder(order.CustomerId, order.MitraId,order,orderDetailMapping)

		return orderResponse, nil
}

func (s *orderService) UpdateStatusOrder(ctx context.Context, orderId primitive.ObjectID, input inputs.UpdateStatusOrderInput) (*mongo.UpdateResult, error){
	
	// status input validation
	if(input.Status != constants.OrderCanceled && input.Status != constants.OrderCompleted && input.Status != constants.OrderProcess && input.Status != constants.OrderRejected){
		return nil, errors.New("status order hanya boleh di update ke antara 'selesai','proses','ditolak','dibatalkan'")
	}

	order, err := s.orderRepository.GetById(ctx, orderId)

	if err != nil{
		return nil, err
	}

	newOrder := bson.M{
		"status" : input.Status,
		"updatedat" : time.Now(),
	}

	updatedOrder, err := s.orderRepository.UpdateOrder(ctx,order.ID, newOrder)

	if err != nil{
		return nil, err
	}

	return updatedOrder, nil
}

func (s *orderService) GetAllOrderMitra(ctx context.Context, userId primitive.ObjectID, status string)([]formatters.OrderResponse, error){

	var orderResponses []formatters.OrderResponse

	status = strings.ToLower(status)
	if(status != "semua" && status != constants.OrderCanceled && status != constants.OrderCompleted && status != constants.OrderProcess && status != constants.OrderRejected && status != constants.OrderWaiting){
		return orderResponses, errors.New("order hanya dapat ditampilkan 'semua' atau antara 'menunggu', 'selesai', 'proses', 'ditolak', 'dibatalkan'")
	}
	
	var bidangData formatters.BidangResponse
	var layananData formatters.LayananDetailMitraResponse

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return nil, errors.New("mitra tidak ditemukan")
	}

	var show string
	if(status == "semua"){
		show = ""
	}else{
		show = status 
	}

	orders, err := s.orderRepository.GetAllOrderMitra(ctx, mitra.ID)

	if err != nil{
		return orderResponses, err
	}

	for _, order := range orders{

		if !strings.Contains(order.Status, show){
			continue
		}

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

	sort.Slice(orderResponses, func(i, j int) bool {
		return orderResponses[i].TerakhirDiUpdate.Unix() < orderResponses[j].TerakhirDiUpdate.Unix()
	})

	return orderResponses, nil
}