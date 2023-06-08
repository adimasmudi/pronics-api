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

type KategoriBidangService interface {
	Save(ctx context.Context, input inputs.AddKategoriInput) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Kategori, error)
	
}

type kategoriBidangService struct{
	kategoriBidangRepository repositories.KategoriRepository
}

func NewKategoriBidangService(kategoriBidangRepository repositories.KategoriRepository) *kategoriBidangService{
	return &kategoriBidangService{kategoriBidangRepository}
}

func (s *kategoriBidangService) Save(ctx context.Context, input inputs.AddKategoriInput) (*mongo.InsertOneResult, error){
	newKategori := models.Kategori{
		ID : primitive.NewObjectID(),
		NamaKategori : input.NamaKategori,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	kategoriAdded, err := s.kategoriBidangRepository.SaveKategori(ctx, newKategori)

	if err != nil{
		return nil, err
	}

	return kategoriAdded, nil
}

func (s *kategoriBidangService) FindAll(ctx context.Context) ([]models.Kategori, error){
	allKategori, err := s.kategoriBidangRepository.FindAll(ctx)

	if err != nil{
		return allKategori, err
	}

	return allKategori, nil
}