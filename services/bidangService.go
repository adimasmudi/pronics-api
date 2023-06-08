package services

import (
	"context"
	"fmt"
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
	FindAll(ctx context.Context) ([]models.Bidang, error)
	
}

type bidangService struct{
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
}

func NewbidangService(bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository) *bidangService{
	return &bidangService{bidangRepository, kategoriRepository}
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
		for _, bidang := range kategoriBidang.BidangId{
			bidangArr = append(bidangArr, bidang)
		}
	}

	bidangArr = append(bidangArr, bidangAdded.InsertedID.(primitive.ObjectID))

	fmt.Println("arr",bidangArr)

	newBidangInKategori := bson.M{
		"bidang_id" : bidangArr,
		"updatedat" : time.Now(),
	}

	insertedBidang, err := s.kategoriRepository.AddBidangToKategori(ctx, kategoriBidang.ID, newBidangInKategori)

	fmt.Println(insertedBidang)

	if err != nil{
		return bidangAdded, err
	}

	return bidangAdded, nil
}

func (s *bidangService) FindAll(ctx context.Context) ([]models.Bidang, error){
	allBidang, err := s.bidangRepository.FindAll(ctx)

	if err != nil{
		return allBidang, err
	}

	return allBidang, nil
}