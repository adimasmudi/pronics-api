package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type WilayahMitraRepository interface {
	SaveWilayahMitra(ctx context.Context, wilayahMitra models.WilayahMitra) (*mongo.InsertOneResult, error)
}

type wilayahMitraRepository struct{
	DB *mongo.Collection
}

func NewWilayahMitraRepository(DB *mongo.Collection) *wilayahMitraRepository{
	return &wilayahMitraRepository{DB}
}

func (r *wilayahMitraRepository) SaveWilayahMitra(ctx context.Context, wilayahMitra models.WilayahMitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, wilayahMitra)

	if err != nil {
		return result, err
	}

	return result, nil
}