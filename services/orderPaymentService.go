package services

import (
	"context"
	"encoding/json"
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
}

type orderPaymentService struct{
	mitraRepository repositories.MitraRepository
	orderRepository repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	orderPaymentRepository repositories.OrderPaymentRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewOrderPaymentService(mitraRepository repositories.MitraRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, orderPaymentRepository repositories.OrderPaymentRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *orderPaymentService{
	return &orderPaymentService{mitraRepository,orderRepository, orderDetailRepository,orderPaymentRepository,bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
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

	mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

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

		biayaLayanan = layananMitra.Harga
	}else{
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

	var transportFee float64

	if distance <= 10{
		transportFee = float64(distance * constants.CostPerKMLessThan10KM)
	}else{
		transportFee = float64(distance * constants.CostPerKMMoreThan10KM)
	}

	var percentageAppsCharge int

	if (biayaLayanan+transportFee) <= 100000{
		percentageAppsCharge = constants.AppsChargePercentageLessThan100k
	}else{
		percentageAppsCharge = constants.AppsChargePercentageMoreThan100k
	}

	biayaAplikasi := (biayaLayanan+transportFee) * float64(percentageAppsCharge) / 100

	if(input.JenisOrder == constants.OrderTakeDelivery){
		transportFee *= 2
	}

	orderPayment, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetailId)

	if err == nil{
		// update
		newOrderPayment := bson.M{
			"biayaperjalanan" : transportFee,
			"biayapelayanan" : biayaLayanan,
			"biayaaplikasi" : biayaAplikasi,
			"totalBiaya" : transportFee + biayaLayanan + biayaAplikasi,
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
	}else{
		layananData.ID = layananToOrder.ID
		layananData.NamaLayanan = layananToOrder.NamaLayanan
		layananData.Harga = layananToOrder.Harga
	}

	orderPaymentToDisplay, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetailId)

	if err != nil{
		return orderResponse, err
	}

	orderPaymentData := helper.MapperOrderPayment(orderPaymentToDisplay)
	orderDetailData := helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPaymentData)
	orderData :=helper.MapperOrder(orderToDisplay.CustomerId, orderToDisplay.MitraId, orderToDisplay, orderDetailData)

	return orderData, nil
}