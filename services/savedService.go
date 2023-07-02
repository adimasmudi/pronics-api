package services

import (
	"context"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SavedService interface {
	Save(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (*mongo.InsertOneResult, error)
	DeleteSaved(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type savedService struct{
	userRepository     repositories.UserRepository
	customerRepository repositories.CustomerRepository
	mitraRepository repositories.MitraRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
	savedRepository repositories.SavedRepository
}

func NewSavedService(userRepository repositories.UserRepository,customerRepository repositories.CustomerRepository, mitraRepository repositories.MitraRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository, savedRepository repositories.SavedRepository) *savedService{
	return &savedService{userRepository, customerRepository,mitraRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, savedRepository}
}

func (s *savedService) Save(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (*mongo.InsertOneResult, error){
	
	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}

	newSaved := models.Saved{
		ID : primitive.NewObjectID(),
		Customer_id: customer.ID,
		Mitra_id: mitraId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	savedAdded, err := s.savedRepository.Save(ctx, newSaved)

	if err != nil{
		return nil, err
	}

	return savedAdded, nil
}

func (s *savedService) DeleteSaved(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error){
	deletedSaved, err := s.savedRepository.Delete(ctx, savedId)

	if err != nil{
		return deletedSaved, err
	}

	return deletedSaved, nil
}