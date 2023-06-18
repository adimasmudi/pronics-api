package services

import (
	"context"
	"pronics-api/inputs"
	"pronics-api/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LayananMitraService interface {
	Save(ctx context.Context, input inputs.AddLayananInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error)
}

type layananMitraService struct{
	layananMitraRepository repositories.LayananMitraRepository
	bidangRepository repositories.BidangRepository
	mitraRepository repositories.MitraRepository
}

func NewLayananMitraService(layananMitraRepository repositories.LayananMitraRepository, bidangRepository repositories.BidangRepository, mitraRepository repositories.MitraRepository) *layananMitraService{
	return &layananMitraService{layananMitraRepository, bidangRepository, mitraRepository}
}

// func (s *layananMitraService) Save(ctx context.Context, input inputs.AddLayananInput, userId primitive.ObjectID) (*mongo.InsertOneResult, error){
// 	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

// 	if err != nil{
// 		return nil, erros.New("Mitra witn the given id not found")
// 	}

// 	newLayananMitra := models.LayananMitra{
// 		ID : primitive.NewObjectID(),
// 		NamaLayanan : input.NamaLayanan,
// 		BidangId: input.BidangId,
// 		MitraId : mitra.ID,
// 		Harga : input.Harga,
// 		AvailableTakeDelivery: input.AvailableTakeDelivery,
// 		CreatedAt: time.Now(),
// 		UpdatedAt : time.Now(),
// 	}


// }
