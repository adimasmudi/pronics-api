package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type KategoriRepository interface {
	SaveKategori(ctx context.Context, kategori models.Kategori) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Kategori, error)
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

func (r *kategoriRepository) FindAll(ctx context.Context) ([]models.Kategori, error){
	var categories []models.Kategori

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return categories, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var kategori models.Kategori
        err := currentRes.Decode(&kategori)
        if err != nil {
            return categories, err
        }

        categories =append(categories, kategori)
    }

	if err := currentRes.Err(); err != nil {
        return categories, err
    }

	currentRes.Close(ctx)

	return categories, nil
}