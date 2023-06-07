package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type KategoriRepository interface {
	SaveKategori(ctx context.Context, kategori models.Kategori) (*mongo.InsertOneResult, error)
}

type kategoriRepository struct{
	DB *mongo.Collection
}

func NewKategoriRepository(DB *mongo.Collection) *kategoriRepository{
	return &kategoriRepository{DB}
}

func (r *kategoriRepository) SaveKategori(ctx context.Context, kategori models.Kategori) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, kategori)

	if err != nil {
		return result, err
	}

	return result, nil
}