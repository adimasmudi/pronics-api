package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type BidangMitraRepository interface {
	SaveBidangMitra(ctx context.Context, bidangMitra models.BidangMitra) (*mongo.InsertOneResult, error)
}

type bidangMitraRepository struct{
	DB *mongo.Collection
}

func NewBidangMitraRepository(DB *mongo.Collection) *bidangMitraRepository{
	return &bidangMitraRepository{DB}
}

func (r *bidangMitraRepository) SaveBidangMitra(ctx context.Context, bidangMitra models.BidangMitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, bidangMitra)

	if err != nil {
		return result, err
	}

	return result, nil
}