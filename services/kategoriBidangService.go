package services

import (
	"context"
	"pronics-api/formatters"
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
	GetKategoriWithBidang(ctx context.Context) ([]formatters.KategoriWithBidangResponse, error)
}

type kategoriBidangService struct{
	kategoriBidangRepository repositories.KategoriRepository
	bidangRepository repositories.BidangRepository
}

func NewKategoriBidangService(kategoriBidangRepository repositories.KategoriRepository, bidangRepository repositories.BidangRepository) *kategoriBidangService{
	return &kategoriBidangService{kategoriBidangRepository, bidangRepository}
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

func (s *kategoriBidangService) GetKategoriWithBidang(ctx context.Context) ([]formatters.KategoriWithBidangResponse, error){
	var kategoriResponses []formatters.KategoriWithBidangResponse
	var kategoriResponse formatters.KategoriWithBidangResponse
	var bidangResponse formatters.BidangResponse

	allKategori, err := s.kategoriBidangRepository.FindAll(ctx)

	if err != nil{
		return kategoriResponses, err
	}

	for _, kategori := range allKategori{
		kategoriResponse.Bidang = nil
		for _, bidang := range kategori.BidangId{
			
			bidang, err := s.bidangRepository.GetById(ctx,bidang)

			if err != nil{
				return kategoriResponses, err
			}

			bidangResponse.ID = bidang.ID
			bidangResponse.NamaBidang = bidang.NamaBidang


			kategoriResponse.Bidang = append(kategoriResponse.Bidang, bidangResponse)
		}



		kategoriResponse.ID = kategori.ID
		kategoriResponse.NamaKategori = kategori.NamaKategori

		kategoriResponses = append(kategoriResponses, kategoriResponse)
	}

	return kategoriResponses, nil
}