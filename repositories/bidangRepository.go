package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidangRepository interface {
	SaveBidang(ctx context.Context, bidang models.Bidang) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Bidang, error)
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

func (r *bidangRepository) FindAll(ctx context.Context) ([]models.Bidang, error){
	var bidangs []models.Bidang

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return bidangs, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var bidang models.Bidang
        err := currentRes.Decode(&bidang)
        if err != nil {
            return bidangs, err
        }

        bidangs =append(bidangs, bidang)
    }

	if err := currentRes.Err(); err != nil {
        return bidangs, err
    }

	currentRes.Close(ctx)

	return bidangs, nil
}