package services

import (
	"context"
	"pronics-api/formatters"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RekeningService interface {
	GetDetailRekening(ctx context.Context, userId primitive.ObjectID) (formatters.RekeningResponse, error)
	UpdateRekening(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateRekeningInput)(*mongo.UpdateResult, error)
	SaveRekening(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateRekeningInput) (*mongo.InsertOneResult, error)
}

type rekeningService struct{
	rekeningRepository repositories.RekeningRepository
	bankRepository repositories.BankRepository
}

func NewRekeningService(rekeningRepository repositories.RekeningRepository, bankRepository repositories.BankRepository) *rekeningService{
	return &rekeningService{rekeningRepository, bankRepository}
}

func (s *rekeningService) SaveRekening(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateRekeningInput) (*mongo.InsertOneResult, error){
	IdBank, err := primitive.ObjectIDFromHex(input.IdBank)

	if err != nil{
		return nil, err
	}
	newRekening := models.Rekening{
		ID : primitive.NewObjectID(),
		UserId: userId,
		BankId : IdBank, 
		NamaPemilik: input.NamaPemilikRekening,
		NomerRekening: input.NomerRekening,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rekeningAdded, err := s.rekeningRepository.SaveRekening(ctx,newRekening)

	if err != nil{
		return nil, err
	}

	return rekeningAdded, nil
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

func (s *rekeningService) UpdateRekening(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateRekeningInput)(*mongo.UpdateResult, error){
	rekening, err := s.rekeningRepository.GetRekeningByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}
	newRekening := bson.M{
		"namapemilik" : input.NamaPemilikRekening,
		"nomerrekening" : input.NomerRekening,
		"bank_id" : input.IdBank,
	}

	updatedRekening, err := s.rekeningRepository.UpdateRekening(ctx,rekening.ID, newRekening)

	if err != nil{
		return nil, err
	}

	return updatedRekening, nil
}