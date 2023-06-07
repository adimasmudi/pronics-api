package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type BidangRepository interface {
	SaveBidang(ctx context.Context, bidang models.Bidang) (*mongo.InsertOneResult, error)
}

type bidangRepository struct{
	DB *mongo.Collection
}

func NewBidangRepository(DB *mongo.Collection) *bidangRepository{
	return &bidangRepository{DB}
}

func (r *bidangRepository) SaveBidang(ctx context.Context, bidang models.Bidang) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, bidang)

	if err != nil {
		return result, err
	}

	return result, nil
}