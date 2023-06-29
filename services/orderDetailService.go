package services

import (
	"context"
	"fmt"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderDetailService interface {
	AddOrUpdateOrderDetail(ctx context.Context, orderId primitive.ObjectID, input inputs.AddOrUpdateOrderDetailInput) (formatters.OrderResponse, error)
}

type orderDetailService struct{
	orderRepository repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewOrderDetailService( orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *orderDetailService{
	return &orderDetailService{orderRepository, orderDetailRepository,bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
}

func (s *orderDetailService) AddOrUpdateOrderDetail(ctx context.Context, orderId primitive.ObjectID, input inputs.AddOrUpdateOrderDetailInput) (formatters.OrderResponse, error){

	var orderResponse formatters.OrderResponse 

	orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, orderId)

	if err == nil{
		// update
		newOrderDetail := bson.M{
			"bidang_id" : input.BidangId,
			"jenisorder" : orderDetail.JenisOrder, // default
			"merk" : input.Merk,
			"layanan_id" : input.LayananId,
			"deskripsikerusakan" : input.DeskripsiKerusakan,
			"alamatpemesanan" : input.AlamatPesanan,
			"updatedat" : time.Now(),
		}

		updatedOrderDetail, err := s.orderDetailRepository.Update(ctx, orderDetail.ID, newOrderDetail)

		if err != nil{
			return orderResponse, err
		}

		fmt.Println(updatedOrderDetail)
	}else{
		// add
		newOrderDetail := models.OrderDetail{
			ID : primitive.NewObjectID(),
			OrderId: orderId,
			BidangId: input.BidangId,
			JenisOrder: constants.OrderHomeCalling,
			Merk : input.Merk,
			LayananId: input.LayananId,
			DeskripsiKerusakan: input.DeskripsiKerusakan,
			AlamatPemesanan: input.AlamatPesanan,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		addedOrderDetail, err := s.orderDetailRepository.Save(ctx, newOrderDetail)

		if err != nil{
			return orderResponse, err
		}

		fmt.Println(addedOrderDetail)
	}
	

	// get order to return
	orderToDisplay, err := s.orderRepository.GetById(ctx, orderId)
	
	if err != nil{
		return orderResponse, err
	}

	orderDetailToDisplay, err := s.orderDetailRepository.GetByOrderId(ctx,orderId)

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

	var orderPayment formatters.OrderPaymentResponse // sementara

	orderDetailData := helper.MapperOrderDetail(orderDetailToDisplay,bidangData,layananData,orderPayment)
	orderData :=helper.MapperOrder(orderToDisplay.CustomerId, orderToDisplay.MitraId, orderToDisplay, orderDetailData)

	return orderData, nil
}