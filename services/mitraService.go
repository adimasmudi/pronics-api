package services

import (
	"context"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MitraService interface {
	GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error)
}

type mitraService struct {
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
}

func NewMitraService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository) *mitraService{
	return &mitraService{userRepository, mitraRepository}
}

func (s *mitraService) GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error){
	var data formatters.MitraResponse

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, user.ID)

	if err != nil{
		return data, err
	}

	data = helper.MapperMitra(user, mitra)

	return data, nil
}