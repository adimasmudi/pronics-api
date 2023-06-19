package services

import (
	"context"
	"errors"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LayananMitraService interface {
	Save(ctx context.Context, input inputs.AddLayananInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error)
	Delete(ctx context.Context, layananMitraId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type layananMitraService struct{
	layananMitraRepository repositories.LayananMitraRepository
	bidangRepository repositories.BidangRepository
	mitraRepository repositories.MitraRepository
}

func NewLayananMitraService(layananMitraRepository repositories.LayananMitraRepository, bidangRepository repositories.BidangRepository, mitraRepository repositories.MitraRepository) *layananMitraService{
	return &layananMitraService{layananMitraRepository, bidangRepository, mitraRepository}
}

// add layanan mitra
func (s *layananMitraService) Save(ctx context.Context, input inputs.AddLayananInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error){
	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return nil, errors.New("mitra with the given id is not found")
	}

	newLayananMitra := models.LayananMitra{
		ID : primitive.NewObjectID(),
		NamaLayanan : input.NamaLayanan,
		BidangId: input.BidangId,
		MitraId : mitra.ID,
		Harga : input.Harga,
		AvailableTakeDelivery: input.AvailableTakeDelivery,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	addedLayananMitra, err := s.layananMitraRepository.Save(ctx, newLayananMitra)

	if err != nil{
		return nil, err
	}

	return addedLayananMitra, nil
}

// get all layanan mitra by bidang by mitra

// get layanan mitra by id

// delete layanan mitra
func (s *layananMitraService) Delete(ctx context.Context, layananMitraId primitive.ObjectID) (*mongo.DeleteResult, error){
	layananMitra, err := s.layananMitraRepository.GetById(ctx, layananMitraId)

	if err != nil{
		return nil, err
	}

	deletedBidang, err := s.layananMitraRepository.DeleteLayananMitra(ctx, layananMitra.ID)

	if err != nil{
		return nil, err
	}

	return deletedBidang, nil
}
