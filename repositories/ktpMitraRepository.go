package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type KTPMitraRepository interface {
	Save(ctx context.Context, eKTPModel models.KTPMitra) (*mongo.InsertOneResult, error)
}

type ktpMitraRepository struct{
	DB *mongo.Collection
}

func NewKTPMitraRepository(DB *mongo.Collection) *ktpMitraRepository{
	return &ktpMitraRepository{DB}
}

func (r *ktpMitraRepository) Save(ctx context.Context, eKTPModel models.KTPMitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, eKTPModel)

	if err != nil {
		return result, err
	}

	return result, nil
}

