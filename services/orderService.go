package services

import (
	"context"
	"errors"
	"pronics-api/configs"
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
	"golang.org/x/exp/slices"
	"googlemaps.github.io/maps"
)

type OrderService interface {
	CreateTemporaryOrder(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (formatters.OrderResponse, error)
	GetAllOrder(ctx context.Context)([]formatters.OrderResponse, error)
	GetOrderDetail(ctx context.Context, orderId primitive.ObjectID) (formatters.OrderResponse, error)
	UpdateStatusOrder(ctx context.Context, orderId primitive.ObjectID, input inputs.UpdateStatusOrderInput) (*mongo.UpdateResult, error)
	GetAllOrderMitra(ctx context.Context, userId primitive.ObjectID, status string) ([]formatters.OrderResponse, error)
	GetDirection(ctx context.Context, userId primitive.ObjectID, orderId primitive.ObjectID)([]maps.Route, error)
	GetAllOrderHistory(ctx context.Context, userId primitive.ObjectID, filter string) ([]formatters.OrderHistoryResponse, error)
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

	userCustomer, err := s.userRepository.GetUserById(ctx, customer.UserId)

	if err != nil{
		return orderData, err
	}

	userMitra, err := s.userRepository.GetUserById(ctx, mitra.UserId)

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

		jarak, err := helper.DistanceCalculation(orderDetailToDisplay.AlamatPemesanan, mitra.Alamat)

		if err != nil{
			return orderData, err
		}

		orderPaymentData := helper.MapperOrderPayment(orderPaymentToDisplay, jarak)
		orderDetailData = helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPaymentData)
		customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
		mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
		orderData = helper.MapperOrder(customerResponse, mitraResponse, order, orderDetailData)

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
	customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
	mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
	orderData = helper.MapperOrder(customerResponse, mitraResponse, orderAdded, orderDetailData)

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
		customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

		if err != nil{
			return orderResponses, err
		}

		mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

		if err != nil{
			return orderResponses, err
		}

		userCustomer, err := s.userRepository.GetUserById(ctx, customer.UserId)

		if err != nil{
			return orderResponses, err
		}

		userMitra, err := s.userRepository.GetUserById(ctx, mitra.UserId)

		if err != nil{
			return orderResponses, err
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

		jarak, err := helper.DistanceCalculation(orderDetail.AlamatPemesanan, mitra.Alamat)

		if err != nil{
			return orderResponses, err
		}

		orderPaymentMapping := helper.MapperOrderPayment(orderPayment, jarak)
		orderDetailMapping := helper.MapperOrderDetail(orderDetail,bidangData,layananData,orderPaymentMapping)
		customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
		mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
		orderMapping := helper.MapperOrder(customerResponse, mitraResponse,order,orderDetailMapping)

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

	customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

	if err != nil{
		return orderResponse, err
	}

	mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

	if err != nil{
		return orderResponse, err
	}

	userCustomer, err := s.userRepository.GetUserById(ctx, customer.UserId)

	if err != nil{
		return orderResponse, err
	}

	userMitra, err := s.userRepository.GetUserById(ctx, mitra.UserId)

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

	jarak, err := helper.DistanceCalculation(orderDetail.AlamatPemesanan, mitra.Alamat)

	if err != nil{
		return orderResponse, err
	}

	orderPaymentMapping := helper.MapperOrderPayment(orderPayment, jarak)
	orderDetailMapping := helper.MapperOrderDetail(orderDetail,bidangData,layananData,orderPaymentMapping)
	customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
	mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
	orderResponse = helper.MapperOrder(customerResponse, mitraResponse,order,orderDetailMapping)

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
	if(status != "semua" && status != constants.OrderCanceled  && status != constants.OrderProcess  && status != constants.OrderWaiting){
		return orderResponses, errors.New("order hanya dapat ditampilkan 'semua' atau antara 'menunggu', 'proses'")
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

		customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

		if err != nil{
			return orderResponses, err
		}

		mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

		if err != nil{
			return orderResponses, err
		}

		userCustomer, err := s.userRepository.GetUserById(ctx, customer.UserId)

		if err != nil{
			return orderResponses, err
		}

		userMitra, err := s.userRepository.GetUserById(ctx, mitra.UserId)

		if err != nil{
			return orderResponses, err
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

		jarak, err := helper.DistanceCalculation(orderDetail.AlamatPemesanan, mitra.Alamat)

		if err != nil{
			return orderResponses, err
		}

		orderPaymentMapping := helper.MapperOrderPayment(orderPayment, jarak)
		orderDetailMapping := helper.MapperOrderDetail(orderDetail,bidangData,layananData,orderPaymentMapping)
		customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
		mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
		orderMapping := helper.MapperOrder(customerResponse, mitraResponse,order,orderDetailMapping)

		orderResponses = append(orderResponses, orderMapping)
	}

	sort.Slice(orderResponses, func(i, j int) bool {
		return orderResponses[i].TerakhirDiUpdate.Unix() < orderResponses[j].TerakhirDiUpdate.Unix()
	})

	return orderResponses, nil
}

func (s *orderService) GetDirection(ctx context.Context, userId primitive.ObjectID, orderId primitive.ObjectID)([]maps.Route, error){

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}

	order, err := s.orderRepository.GetById(ctx, orderId)

	if err != nil{
		return nil, err
	}

	orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

	if err != nil{
		return nil, err
	}

	c := configs.InitMap()
	// get directions
	r := &maps.DirectionsRequest{
		Origin:      mitra.Alamat,
		Destination: orderDetail.AlamatPemesanan,
	}

	route, _, err := c.Directions(ctx, r)
	if err != nil {
		return route, err
	}

	return route, err
}

func (s *orderService) GetAllOrderHistory(ctx context.Context, userId primitive.ObjectID, filter string) ([]formatters.OrderHistoryResponse, error){
	var allHistories []formatters.OrderHistoryResponse

	if filter == ""{
		return allHistories, errors.New("filter harus diisi dengan nilai antara 'semua','selesai','menunggu','proses','ditolak','dibatalkan'")
	}

	user, err := s.userRepository.GetUserById(ctx, userId)

	if err != nil{
		return allHistories, err
	}

	var typeUser string
	var orders []models.Order
	if user.Type == constants.UserMitra{

		if filter != "semua" && filter != constants.OrderCanceled && filter != constants.OrderRejected && filter != constants.OrderCompleted{
			return allHistories, errors.New("history mitra hanya bisa dengan filter 'semua', 'dibatalkan', 'ditolak', 'selesai'")
		}

		mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

		if err != nil{
			return allHistories, err
		}

		typeUser = user.Type

		orders, err = s.orderRepository.GetAllOrderMitra(ctx, mitra.ID)

		if err != nil{
			return allHistories, err
		}

		for idx,order := range orders{
			if order.Status == constants.OrderProcess || order.Status == constants.OrderWaiting || order.Status == constants.OrderTemporary{
				orders = slices.Delete(orders, idx, idx + 1)
			}
		}

	}else if user.Type == constants.UserCustomer{
		
		typeUser = user.Type

		customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

		if err != nil{
			return allHistories, err
		}

		orders, err = s.orderRepository.GetAllOrderCustomer(ctx, customer.ID)

		if err != nil{
			return allHistories, err
		}

		for idx,order := range orders{
			if order.Status == constants.OrderTemporary{
				orders = slices.Delete(orders, idx, idx + 1)
			}
		}
	}else{
		return allHistories, err
	}
	
	if(filter != "semua"){
		for idx,order := range orders{
			if order.Status != filter{
				orders = slices.Delete(orders, idx, idx + 1)
			}
		}
	}

	for _, order := range orders{
		var historyResponse formatters.OrderHistoryResponse

		var name string
		var namaToko = ""
		var idUser primitive.ObjectID

		if typeUser == constants.UserMitra{
			customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

			if err != nil{
				return allHistories, err
			}

			idUser = customer.UserId
		}else if typeUser == constants.UserCustomer{
			mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

			if err != nil{
				return allHistories, err
			}

			idUser = mitra.UserId
			namaToko = mitra.NamaToko
		}else{
			return allHistories, errors.New("type user harus antara mitra atau customer untuk melakukan aksi ini")
		}

		user, err := s.userRepository.GetUserById(ctx, idUser)

		if err != nil{
			return allHistories, err
		}

		name = user.NamaLengkap
		
		orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

		if err != nil{
			return allHistories, err
		}

		orderPayment, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetail.ID)

		if err != nil{
			return allHistories, err
		}

		layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)
		var namaLayanan string

		if err != nil{
			layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

			if err != nil{
				return allHistories, err
			}

			namaLayanan = layananMitra.NamaLayanan
		}else{
			namaLayanan = layanan.NamaLayanan
		}

		historyResponse.ID = order.ID

		if namaToko != ""{
			historyResponse.Name = namaToko
		}else{
			historyResponse.Name = name
		}

		historyResponse.Status = order.Status
		historyResponse.AlamatPemesanan = orderDetail.AlamatPemesanan
		historyResponse.Layanan = namaLayanan
		historyResponse.TotalBayar = int(orderPayment.TotalBiaya)
		historyResponse.TanggalOrder = order.TanggalOrderSelesai
		
		allHistories = append(allHistories, historyResponse)

	}

	sort.Slice(allHistories, func(i, j int) bool {
		return allHistories[i].TanggalOrder.Unix() < allHistories[j].TanggalOrder.Unix()
	})

	return allHistories, nil
}