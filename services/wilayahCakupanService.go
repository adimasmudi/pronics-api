package services

import (
	"context"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WilayahCakupanService interface {
	Save(ctx context.Context, input inputs.AddWilayahCakupanInput) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.WilayahCakupan, error)
	
}

type wilayahCakupanService struct{
	wilayahCakupanRepository repositories.WilayahRepository
}

func NewWilayahCakupanService(wilayahCakupanRepository repositories.WilayahRepository) *wilayahCakupanService{
	return &wilayahCakupanService{wilayahCakupanRepository}
}

func (s *wilayahCakupanService) Save(ctx context.Context, input inputs.AddWilayahCakupanInput) (*mongo.InsertOneResult, error){
	newWilayah := models.WilayahCakupan{
		ID : primitive.NewObjectID(),
		NamaWilayah : input.NamaWilayah,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	wilayahAdded, err := s.wilayahCakupanRepository.SaveWilayah(ctx, newWilayah)

	if err != nil{
		return nil, err
	}

	return wilayahAdded, nil
}

func (s *wilayahCakupanService) FindAll(ctx context.Context) ([]models.WilayahCakupan, error){
	allWilayah, err := s.wilayahCakupanRepository.FindAll(ctx)

	if err != nil{
		return allWilayah, err
	}

	return allWilayah, nil
}