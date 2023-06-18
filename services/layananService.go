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
	FindAll(ctx context.Context) ([]models.Layanan, error)
	FindById(ctx context.Context, layananId primitive.ObjectID) (models.Layanan, error)
	DeleteLayanan(ctx context.Context, layananId primitive.ObjectID)(*mongo.DeleteResult, error)
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

func (s *layananService) FindAll(ctx context.Context) ([]models.Layanan, error){
	allLayanan, err := s.layananRepository.FindAll(ctx)

	if err != nil{
		return allLayanan, err
	}

	return allLayanan, nil
}

func (s *layananService) FindById(ctx context.Context, layananId primitive.ObjectID) (models.Layanan, error){
	layanan, err := s.layananRepository.GetById(ctx, layananId)

	if err != nil{
		return layanan, err
	}

	return layanan, nil
}

func (s *layananService) DeleteLayanan(ctx context.Context, layananId primitive.ObjectID)(*mongo.DeleteResult, error){
	layanan, err := s.layananRepository.GetById(ctx, layananId)

	if err != nil{
		return nil, err
	}

	bidang, err := s.bidangRepository.GetById(ctx, layanan.BidangId)

	if err != nil{
		return nil, err
	}

	oldLayananInBidangArr := []primitive.ObjectID{}

	for _, item := range bidang.LayananId{
		if item != layanan.ID{
			oldLayananInBidangArr = append(oldLayananInBidangArr, item)
		}
	}

	var oldBidangUpdate primitive.M

	oldBidangUpdate = bson.M{
		"layanan_id" : oldLayananInBidangArr,
		"updatedat" : time.Now(),
	}

	oldUpdated, err := s.bidangRepository.UpdateBidang(ctx, bidang.ID,oldBidangUpdate)

	if err != nil{
		return nil, err
	}

	fmt.Println(oldUpdated)

	deletedLayanan, err := s.layananRepository.DeleteLayanan(ctx, layananId)

	if err != nil{
		return nil, err
	}

	return deletedLayanan, nil
}
