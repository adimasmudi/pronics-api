package services

import (
	"context"
	"os"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BankService interface {
	SaveBank(ctx context.Context, input inputs.AddBankInput, fileName string) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Bank, error)
}

type bankService struct{
	bankRepository repositories.BankRepository
}

func NewBankService(bankRepository repositories.BankRepository) *bankService{
	return &bankService{bankRepository}
}

func (s *bankService) SaveBank(ctx context.Context, input inputs.AddBankInput, fileName string) (*mongo.InsertOneResult, error){
	var newBank models.Bank
	if fileName != ""{
		newBank = models.Bank{
			ID : primitive.NewObjectID(),
			NamaBank : input.NamaBank,
			LogoBank: os.Getenv("CLOUD_STORAGE_READ_LINK")+"bank/"+fileName,
		}
	}else{
		newBank = models.Bank{
			ID : primitive.NewObjectID(),
			NamaBank : input.NamaBank,
		}
	}

	newBank.CreatedAt = time.Now()
	newBank.UpdatedAt = time.Now()

	addedBank, err := s.bankRepository.Save(ctx, newBank)

	if err != nil{
		return addedBank, err
	}

	return addedBank, nil
}

func (s *bankService) FindAll(ctx context.Context) ([]models.Bank, error){
	allBank, err := s.bankRepository.FindAll(ctx)

	if err != nil{
		return allBank, err
	}

	return allBank, nil
}