package services

import (
	"context"
	"pronics-api/formatters"
	"pronics-api/repositories"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RekeningService interface {
	GetDetailRekening(ctx context.Context, userId primitive.ObjectID) (formatters.RekeningResponse, error)
}

type rekeningService struct{
	rekeningRepository repositories.RekeningRepository
	bankRepository repositories.BankRepository
}

func NewRekeningService(rekeningRepository repositories.RekeningRepository, bankRepository repositories.BankRepository) *rekeningService{
	return &rekeningService{rekeningRepository, bankRepository}
}

func (s *rekeningService) GetDetailRekening(ctx context.Context, userId primitive.ObjectID) (formatters.RekeningResponse, error){
	var data formatters.RekeningResponse

	rekening, err := s.rekeningRepository.GetRekeningByIdUser(ctx, userId)

	if err != nil{
		return data, err
	}

	bank, err := s.bankRepository.GetBankById(ctx, rekening.BankId)

	if err != nil{
		return data, err
	}

	data.ID = rekening.ID
	data.Bank = bank
	data.NamaPemilik = rekening.NamaPemilik
	data.NomerRekening = rekening.NomerRekening

	return data, nil
}