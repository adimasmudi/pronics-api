package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type BankRepository interface {
	Save(ctx context.Context, bank models.Bank) (*mongo.InsertOneResult, error)
}

type bankRepository struct{
	DB *mongo.Collection
}

func NewBankRepository(DB *mongo.Collection) *bankRepository{
	return &bankRepository{DB}
}

func (r *bankRepository) Save(ctx context.Context, bank models.Bank) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, bank)

	if err != nil {
		return result, err
	}

	return result, nil
}