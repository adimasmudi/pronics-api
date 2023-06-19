package services

import (
	"context"
	"fmt"
	"pronics-api/formatters"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidangService interface {
	SaveBidang(ctx context.Context, input inputs.AddBidangInput, creator_id primitive.ObjectID) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]formatters.BidangResponse, error)
	FindById(ctx context.Context, bidangId primitive.ObjectID) (formatters.DetailBidangResponse, error)
	UpdateBidang(ctx context.Context, editor_id primitive.ObjectID, bidangId primitive.ObjectID, input inputs.AddBidangInput) (*mongo.UpdateResult, error)
	DeleteBidang(ctx context.Context, bidangId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type bidangService struct{
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
}

func NewbidangService(bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository) *bidangService{
	return &bidangService{bidangRepository, kategoriRepository, layananRepository}
}

func (s *bidangService) SaveBidang(ctx context.Context, input inputs.AddBidangInput, creator_id primitive.ObjectID) (*mongo.InsertOneResult, error){
	newBidang := models.Bidang{
		ID : primitive.NewObjectID(),
		NamaBidang : input.NamaBidang,
		KategoriId: input.KategoriId,
		CreatedById: creator_id,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	kategoriBidang, err := s.kategoriRepository.GetById(ctx, input.KategoriId)

	if err != nil{
		return nil, err
	}

	bidangAdded, err := s.bidangRepository.SaveBidang(ctx, newBidang)

	if err != nil{
		return nil, err
	}

	var bidangArr []primitive.ObjectID

	if kategoriBidang.BidangId != nil{
		bidangArr = append(bidangArr, kategoriBidang.BidangId...)
	}

	bidangArr = append(bidangArr, bidangAdded.InsertedID.(primitive.ObjectID))

	newBidangInKategori := bson.M{
		"bidang_id" : bidangArr,
		"updatedat" : time.Now(),
	}

	insertedBidang, err := s.kategoriRepository.AddBidangToKategori(ctx, kategoriBidang.ID, newBidangInKategori)

	if err != nil{
		return bidangAdded, err
	}

	fmt.Println(insertedBidang)

	return bidangAdded, nil
}

func (s *bidangService) FindAll(ctx context.Context) ([]formatters.BidangResponse, error){
	var BidangResponses []formatters.BidangResponse
	var BidangResponse formatters.BidangResponse

	allBidang, err := s.bidangRepository.FindAll(ctx)

	if err != nil{
		return BidangResponses, err
	}

	for _, bidang := range allBidang{
		kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

		if err != nil{
			return BidangResponses, err
		}


		BidangResponse.ID = bidang.ID
		BidangResponse.NamaBidang = bidang.NamaBidang
		BidangResponse.Kategori = kategori.NamaKategori

		BidangResponses = append(BidangResponses, BidangResponse)
	}

	return BidangResponses, nil
}

func (s *bidangService) FindById(ctx context.Context, bidangId primitive.ObjectID) (formatters.DetailBidangResponse, error){
	var detailBidang formatters.DetailBidangResponse
	bidang, err := s.bidangRepository.GetById(ctx, bidangId)

	if err != nil{
		return detailBidang, err
	}

	kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

	if err != nil{
		return detailBidang, err
	}

	layananInBidang, err := s.layananRepository.FindAllByBidangId(ctx, bidang.ID)

	if err != nil{
		return detailBidang, err
	}

	var allLayanan []formatters.LayananResponse

	for _, layananItem := range layananInBidang{
		var layanan formatters.LayananResponse

		layanan.ID = layananItem.ID
		layanan.NamaLayanan = layananItem.NamaLayanan

		allLayanan = append(allLayanan, layanan)
	}

	detailBidang.ID = bidang.ID
	detailBidang.NamaBidang = bidang.NamaBidang
	detailBidang.Kategori = kategori.NamaKategori
	detailBidang.Layanan = allLayanan

	return detailBidang, nil
}


func (s *bidangService) UpdateBidang(ctx context.Context, editor_id primitive.ObjectID, bidangId primitive.ObjectID, input inputs.AddBidangInput) (*mongo.UpdateResult, error){
	var newBidang primitive.M

	currentBidang, err := s.bidangRepository.GetById(ctx, bidangId)

	if err != nil{
		return nil, err
	}

	newBidang = bson.M{
		"namabidang" : input.NamaBidang,
		"kategori_id" : input.KategoriId,
		"createdbyid" : editor_id,
		"updatedat" : time.Now(),
	}

	if currentBidang.KategoriId != input.KategoriId{
		// update old kategori
		oldKategori, err := s.kategoriRepository.GetById(ctx, currentBidang.KategoriId)
		oldBidangInKategoriArr := []primitive.ObjectID{}

		if err != nil{
			return nil, err
		}

		for _, item := range oldKategori.BidangId{
			if item != currentBidang.ID{
				oldBidangInKategoriArr = append(oldBidangInKategoriArr, item)
			}
		}

		var oldKategoriUpdate primitive.M

		oldKategoriUpdate = bson.M{
			"bidang_id" : oldBidangInKategoriArr,
			"updatedat" : time.Now(),
		}

		oldUpdated, err := s.kategoriRepository.AddBidangToKategori(ctx, oldKategori.ID,oldKategoriUpdate)

		if err != nil{
			return nil, err
		}

		fmt.Println(oldUpdated)

		// ambil kategori baru
		newKategori, err := s.kategoriRepository.GetById(ctx, input.KategoriId)

		if err != nil{
			return nil, err
		}

		newKategoriBidangArr := []primitive.ObjectID{}
		newKategoriBidangArr = append(newKategoriBidangArr, newKategori.BidangId...)
		newKategoriBidangArr = append(newKategoriBidangArr, currentBidang.ID)

		// update dengan masukkan id bidang yang sekarang ke kategori lama

		var newKategoriUpdate primitive.M

		newKategoriUpdate = bson.M{
			"bidang_id" : newKategoriBidangArr,
			"updatedat" : time.Now(),
		}

		newUpdated, err := s.kategoriRepository.AddBidangToKategori(ctx, newKategori.ID,newKategoriUpdate)

		if err != nil{
			return nil, err
		}

		fmt.Println(newUpdated)
	}

	updatedBidang, err := s.bidangRepository.UpdateBidang(ctx, bidangId, newBidang)

	if err != nil{
		return nil, err
	}

	return updatedBidang, nil
}

func (s *bidangService) DeleteBidang(ctx context.Context, bidangId primitive.ObjectID) (*mongo.DeleteResult, error){
	
	bidang, err := s.bidangRepository.GetById(ctx, bidangId)

	if err != nil{
		return nil, err
	}

	kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

	if err != nil{
		return nil, err
	}

	oldBidangInKategoriArr := []primitive.ObjectID{}

	for _, item := range kategori.BidangId{
		if item != bidang.ID{
			oldBidangInKategoriArr = append(oldBidangInKategoriArr, item)
		}
	}

	var oldKategoriUpdate primitive.M

	oldKategoriUpdate = bson.M{
		"bidang_id" : oldBidangInKategoriArr,
		"updatedat" : time.Now(),
	}

	oldUpdated, err := s.kategoriRepository.AddBidangToKategori(ctx, kategori.ID,oldKategoriUpdate)

	if err != nil{
		return nil, err
	}

	fmt.Println(oldUpdated)


	deletedBidang, err := s.bidangRepository.DeleteBidang(ctx, bidangId)

	if err != nil{
		return nil, err
	}

	return deletedBidang, nil
}