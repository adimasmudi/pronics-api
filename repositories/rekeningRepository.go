package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type RekeningRepository interface {
	SaveRekening(ctx context.Context, rekening models.Rekening) (*mongo.InsertOneResult, error)
}

type rekeningRepository struct{
	DB *mongo.Collection
}

func NewRekeningRepository(DB *mongo.Collection) *rekeningRepository{
	return &rekeningRepository{DB}
}

func (r *rekeningRepository) SaveRekening(ctx context.Context, rekening models.Rekening) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, rekening)

	if err != nil {
		return result, err
	}

	return result, nil
}