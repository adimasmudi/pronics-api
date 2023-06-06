package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type MitraRepository interface {
	SaveMitra(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error)
}

type mitraRepository struct{
	DB *mongo.Collection
}

func NewMitraRepository(DB *mongo.Collection) *mitraRepository{
	return &mitraRepository{DB}
}

func (r *mitraRepository) SaveMitra(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, mitra)

	if err != nil {
		return result, err
	}

	return result, nil
}