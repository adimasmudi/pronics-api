package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderPaymentService interface {
	AddOrUpdateOrderPayment(ctx context.Context, orderDetailId primitive.ObjectID, input inputs.AddOrUpdateOrderPaymentInput) (formatters.OrderResponse, error)
	ConfirmPayment(ctx context.Context, orderPaymentId primitive.ObjectID,input inputs.ConfirmPaymentInput, fileName string) (formatters.OrderResponse, error)
}

type orderPaymentService struct{
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

func NewOrderPaymentService(userRepository repositories.UserRepository,mitraRepository repositories.MitraRepository,customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, orderPaymentRepository repositories.OrderPaymentRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *orderPaymentService{
	return &orderPaymentService{userRepository,mitraRepository,customerRepository,orderRepository, orderDetailRepository,orderPaymentRepository,bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
}

func (s *orderPaymentService) AddOrUpdateOrderPayment(ctx context.Context, orderDetailId primitive.ObjectID, input inputs.AddOrUpdateOrderPaymentInput) (formatters.OrderResponse, error){

	var orderResponse formatters.OrderResponse

	orderDetail, err := s.orderDetailRepository.GetById(ctx, orderDetailId)

	if err != nil{
		return orderResponse, err
	}

	order, err := s.orderRepository.GetById(ctx, orderDetail.OrderId)

	if err != nil{
		return orderResponse, err
	}

	if order.Status != constants.OrderTemporary{
		return orderResponse, errors.New("order sudah diproses, bukan temporary order")
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

	layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)
	var biayaLayanan float64

	if err != nil{
		layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

		if err != nil{
			return orderResponse, err
		}

		if(input.JenisOrder == constants.OrderTakeDelivery && !layananMitra.AvailableTakeDelivery){
			return orderResponse, errors.New("order ini tidak menerima take & delivery")
		}

		biayaLayanan = layananMitra.Harga
	}else{
		if(input.JenisOrder == constants.OrderTakeDelivery && !layanan.AvailableTakeDelivery){
			return orderResponse, errors.New("order ini tidak menerima take & delivery")
		}
		biayaLayanan = layanan.Harga
	}

	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&units=imperial&key=%s",mitra.Alamat,orderDetail.AlamatPemesanan,os.Getenv("MAPS_API_KEY"))
	url = strings.Replace(url, " ", "%20", -1)
	method := "GET"


	client := &http.Client {}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return orderResponse, err
	}
	res, err := client.Do(req)
	if err != nil {
		return orderResponse, err
	}
	
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return orderResponse, err
	}

	var distanceMatrix helper.DistanceMatrixResult
	
	err = json.Unmarshal(body, &distanceMatrix)
	if err != nil {
		return orderResponse, err
	}

	distance := distanceMatrix.Rows[0].Elements[0].Distance.Value / 1000

	if (input.JenisOrder != constants.OrderTakeDelivery && input.JenisOrder != constants.OrderHomeCalling){
		return orderResponse, errors.New("jenis order harus antara home calling atau take & delivery")
	}

	transportFee := 1.0

	if(input.JenisOrder == constants.OrderTakeDelivery){
		transportFee *= 2
	}

	if(input.JenisOrder != orderDetail.JenisOrder){
		newOrderDetail := bson.M{
			"jenisorder" : input.JenisOrder,
			"updatedat" : time.Now(),
		}
	
		updatedOrderDetail, err := s.orderDetailRepository.Update(ctx, orderDetailId, newOrderDetail)
	
		if err != nil{
			return orderResponse, err
		}
	
		fmt.Println(updatedOrderDetail)
	}

	if distance <= 10{
		transportFee *= float64(distance * constants.CostPerKMLessThan10KM)
	}else{
		transportFee *= float64(distance * constants.CostPerKMMoreThan10KM)
	}

	var percentageAppsCharge int

	fmt.Println("ttl",biayaLayanan+transportFee)

	if (biayaLayanan+transportFee) <= 100000{
		percentageAppsCharge = constants.AppsChargePercentageLessThan100k
	}else{
		percentageAppsCharge = constants.AppsChargePercentageMoreThan100k
	}

	fmt.Println(percentageAppsCharge)

	biayaAplikasi := (biayaLayanan+transportFee) * float64(percentageAppsCharge) / 100

	orderPayment, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetailId)

	if err == nil{
		// update
		newOrderPayment := bson.M{
			"biayaperjalanan" : transportFee,
			"biayapelayanan" : biayaLayanan,
			"biayaaplikasi" : biayaAplikasi,
			"totalbiaya" : transportFee + biayaLayanan + biayaAplikasi,
			"updatedat" : time.Now(),
		}

		updatedOrderPayment, err := s.orderPaymentRepository.Update(ctx, orderPayment.ID, newOrderPayment)

		if err != nil{
			return orderResponse, err
		}

		fmt.Println(updatedOrderPayment)

	}else{
		// add
		newOrderPayment := models.OrderPayment{
			ID : primitive.NewObjectID(),
			OrderDetailId: orderDetailId,
			BiayaPelayanan: biayaLayanan,
			BiayaPerjalanan: transportFee,
			BiayaAplikasi: biayaAplikasi,
			TotalBiaya: transportFee + biayaLayanan + biayaAplikasi,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		addedOrderPayment, err := s.orderPaymentRepository.Save(ctx, newOrderPayment)

		if err != nil{
			return orderResponse, err
		}

		fmt.Println(addedOrderPayment)
	}
	

	// get order to return
	orderToDisplay, err := s.orderRepository.GetById(ctx, order.ID)
	
	if err != nil{
		return orderResponse, err
	}

	orderDetailToDisplay, err := s.orderDetailRepository.GetByOrderId(ctx,order.ID)

	if err != nil{
		return orderResponse, err
	}

	var bidangData formatters.BidangResponse

	bidangToOrder, err := s.bidangRepository.GetById(ctx, orderDetailToDisplay.BidangId)

	if err != nil{
		return orderResponse, err
	}

	kategoriToOrder, err := s.kategoriRepository.GetById(ctx, bidangToOrder.KategoriId)

	if err != nil{
		return orderResponse, err
	}

	bidangData.ID = bidangToOrder.ID
	bidangData.Kategori = kategoriToOrder.NamaKategori
	bidangData.NamaBidang = bidangToOrder.NamaBidang

	var layananData formatters.LayananDetailMitraResponse

	layananToOrder, err := s.layananRepository.GetById(ctx, orderDetailToDisplay.LayananId)

	if err != nil{
		layananMitraToOrder, err := s.layananMitraRepository.GetById(ctx, orderDetailToDisplay.LayananId)

		if err != nil{
			return orderResponse, err
		}

		layananData.ID = layananMitraToOrder.ID
		layananData.NamaLayanan = layananMitraToOrder.NamaLayanan
		layananData.Harga = layananMitraToOrder.Harga
		layananData.BidangId = bidangToOrder.ID
	}else{
		layananData.ID = layananToOrder.ID
		layananData.NamaLayanan = layananToOrder.NamaLayanan
		layananData.Harga = layananToOrder.Harga
		layananData.BidangId = bidangToOrder.ID
	}

	orderPaymentToDisplay, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetailId)

	if err != nil{
		return orderResponse, err
	}

	jarak, err := helper.DistanceCalculation(orderDetailToDisplay.AlamatPemesanan,mitra.Alamat)

	if err != nil{
		return orderResponse, err
	}

	orderPaymentData := helper.MapperOrderPayment(orderPaymentToDisplay, jarak)
	orderDetailData := helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPaymentData)
	customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
	mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
	orderData :=helper.MapperOrder(customerResponse, mitraResponse, orderToDisplay, orderDetailData)

	return orderData, nil
}

func (s *orderPaymentService) ConfirmPayment(ctx context.Context, orderPaymentId primitive.ObjectID,input inputs.ConfirmPaymentInput, fileName string) (formatters.OrderResponse, error){
	var orderData formatters.OrderResponse
	var buktiBayar string

	orderPayment, err := s.orderPaymentRepository.GetById(ctx, orderPaymentId)

	if err != nil{
		return orderData, err
	}

	orderDetail, err := s.orderDetailRepository.GetById(ctx, orderPayment.OrderDetailId)

	if err != nil{
		return orderData, err
	}

	order, err := s.orderRepository.GetById(ctx, orderDetail.OrderId)

	if err != nil{
		return orderData, err
	}

	if order.Status != constants.OrderTemporary{
		return orderData, errors.New("order sudah diproses, bukan temporary order")
	}

	customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

	if err != nil{
		return orderData, err
	}

	mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

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

	if fileName != ""{
		buktiBayar = os.Getenv("CLOUD_STORAGE_READ_LINK")+"buktiBayar/"+fileName
	}

	if(input.MetodePembayaran != constants.AutomaticPayment && input.MetodePembayaran != constants.BankTransferPayment && input.MetodePembayaran != constants.CashPayment){
		return orderData, errors.New("metode pembayaran harus antara 'bayar otomatis','bank transfer','cash' (lowercase)")
	}

	newOrderPayment := bson.M{
		"metodepembayaran" : input.MetodePembayaran,
		"buktibayar" : buktiBayar,
		"updatedat" : time.Now(),
	}

	updatedPayment, err := s.orderPaymentRepository.Update(ctx, orderPaymentId, newOrderPayment)

	if err != nil{
		return orderData, err
	}

	fmt.Println(updatedPayment)

	transaksiId := fmt.Sprintf("PRN-%s%s",string(order.ID.Hex())[0:len(string(order.ID.Hex()))/2],string(time.Now().UTC().Format("2023-06-30")))

	newOrder := bson.M{
		"status" : constants.OrderWaiting,
		"transaksiid" : transaksiId,
		"tanggalorderselesai" : time.Now(),
		"updatedat" : time.Now(),
	}

	updatedOrder, err := s.orderRepository.UpdateOrder(ctx,order.ID, newOrder)

	if err != nil{
		return orderData, err
	}

	fmt.Println(updatedOrder)

	// get order to return
	orderToDisplay, err := s.orderRepository.GetById(ctx, order.ID)
	
	if err != nil{
		return orderData, err
	}

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

	var layananData formatters.LayananDetailMitraResponse

	layananToOrder, err := s.layananRepository.GetById(ctx, orderDetailToDisplay.LayananId)

	if err != nil{
		layananMitraToOrder, err := s.layananMitraRepository.GetById(ctx, orderDetailToDisplay.LayananId)

		if err != nil{
			return orderData, err
		}

		layananData.ID = layananMitraToOrder.ID
		layananData.NamaLayanan = layananMitraToOrder.NamaLayanan
		layananData.Harga = layananMitraToOrder.Harga
		layananData.BidangId = bidangToOrder.ID
	}else{
		layananData.ID = layananToOrder.ID
		layananData.NamaLayanan = layananToOrder.NamaLayanan
		layananData.Harga = layananToOrder.Harga
		layananData.BidangId = bidangToOrder.ID
	}

	orderPaymentToDisplay, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetail.ID)

	if err != nil{
		return orderData, err
	}

	jarak, err := helper.DistanceCalculation(orderDetailToDisplay.AlamatPemesanan, mitra.Alamat)

	if err != nil{
		return orderData, err
	}
	
	orderPaymentData := helper.MapperOrderPayment(orderPaymentToDisplay, jarak)
	orderDetailData := helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPaymentData)
	customerResponse := helper.MapperCustomer(userCustomer,customer,nil)
	mitraResponse := helper.MapperMitra(userMitra, mitra,models.WilayahCakupan{}, nil)
	orderData =helper.MapperOrder(customerResponse, mitraResponse, orderToDisplay, orderDetailData)

	return orderData, nil
}