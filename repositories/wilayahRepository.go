package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type WilayahRepository interface {
	SaveWilayah(ctx context.Context, wilayah models.WilayahCakupan) (*mongo.InsertOneResult, error)
}

type wilayahRepository struct{
	DB *mongo.Collection
}

func NewWilayahRepository(DB *mongo.Collection) *wilayahRepository{
	return &wilayahRepository{DB}
}

func (r *wilayahRepository) SaveWilayah(ctx context.Context, wilayah models.WilayahCakupan) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, wilayah)

	if err != nil {
		return result, err
	}

	return result, nil
}