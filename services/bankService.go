package services

import (
	"context"
	"os"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BankService interface {
	SaveBank(ctx context.Context, input inputs.AddBankInput, fileName string) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Bank, error)
	GetDetail(ctx context.Context, bankId primitive.ObjectID) (models.Bank, error)
	UpdateBank(ctx context.Context, bankId primitive.ObjectID, input inputs.AddBankInput, fileName string) (*mongo.UpdateResult, error)
	DeleteBank(ctx context.Context, bankId primitive.ObjectID) (*mongo.DeleteResult, error)
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

func (s *bankService) GetDetail(ctx context.Context, bankId primitive.ObjectID) (models.Bank, error){
	bank, err := s.bankRepository.GetBankById(ctx, bankId)

	if err != nil{
		return bank, err
	}

	return bank, nil
}

func (s *bankService) UpdateBank(ctx context.Context, bankId primitive.ObjectID, input inputs.AddBankInput, fileName string) (*mongo.UpdateResult, error){
	var newBank primitive.M

	bank, err := s.bankRepository.GetBankById(ctx, bankId)

	if err != nil{
		return nil, err
	}

	if fileName == ""{
		fileName = bank.LogoBank
	}else{
		fileName = os.Getenv("CLOUD_STORAGE_READ_LINK")+"bank/"+fileName
	}

	newBank = bson.M{
		"namabank" : input.NamaBank,
		"logobank" : fileName,
		"updatedat" : time.Now(),
	}

	updatedBank, err := s.bankRepository.UpdateBank(ctx, bankId,newBank)

	if err != nil{
		return nil, err
	}

	return updatedBank, nil
}

func (s *bankService) DeleteBank(ctx context.Context, bankId primitive.ObjectID) (*mongo.DeleteResult, error){
	deletedBank, err := s.bankRepository.DeleteBank(ctx, bankId)

	if err != nil{
		return deletedBank, err
	}

	return deletedBank, nil
}