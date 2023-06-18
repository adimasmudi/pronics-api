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

type LayananService interface {
	SaveLayanan(ctx context.Context, input inputs.AddLayananInput) (*mongo.InsertOneResult, error)
	
}

type layananService struct{
	layananRepository repositories.LayananRepository
	bidangRepository repositories.BidangRepository
}

func NewLayananService(layananRepository repositories.LayananRepository, bidangRepository repositories.BidangRepository) *layananService{
	return &layananService{layananRepository, bidangRepository}
}

func (s *layananService) SaveLayanan(ctx context.Context, input inputs.AddLayananInput) (*mongo.InsertOneResult, error){
	
	newLayanan := models.Layanan{
		ID : primitive.NewObjectID(),
		NamaLayanan : input.NamaLayanan,
		BidangId: input.BidangId,
		Harga : input.Harga,
		AvailableTakeDelivery: input.AvailableTakeDelivery,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	bidang, err := s.bidangRepository.GetById(ctx, input.BidangId)

	if err != nil{
		return nil, err
	}

	layananAdded, err := s.layananRepository.SaveLayanan(ctx, newLayanan)

	if err != nil{
		return nil, err
	}

	var layananArr []primitive.ObjectID

	if bidang.LayananId != nil{
		layananArr = append(layananArr, bidang.LayananId...)
	}

	layananArr = append(layananArr, layananAdded.InsertedID.(primitive.ObjectID))

	newLayananInBidang := bson.M{
		"layanan_id" : layananArr,
		"updatedat" : time.Now(),
	}

	insertedLayanan, err := s.bidangRepository.UpdateBidang(ctx, bidang.ID,newLayananInBidang)

	if err != nil{
		return nil, err
	}

	fmt.Println(insertedLayanan)
	
	return layananAdded, nil
}
