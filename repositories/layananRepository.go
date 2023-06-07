package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type LayananRepository interface {
	SaveLayanan(ctx context.Context, layanan models.Layanan) (*mongo.InsertOneResult, error)
}

type layananRepository struct{
	DB *mongo.Collection
}

func NewLayananRepository(DB *mongo.Collection) *layananRepository{
	return &layananRepository{DB}
}

func (r *layananRepository) SaveLayanan(ctx context.Context, layanan models.Layanan) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, layanan)

	if err != nil {
		return result, err
	}

	return result, nil
}