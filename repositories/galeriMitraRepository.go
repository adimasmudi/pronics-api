package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type GaleriRepository interface {
	Save(ctx context.Context, galeri models.GaleriMitra) (*mongo.InsertOneResult, error)
	// FindAll(ctx context.Context) ([]models.Kategori, error)
}

type galeriRepository struct {
	DB *mongo.Collection
}

func NewGaleriRepository(DB *mongo.Collection) *galeriRepository {
	return &galeriRepository{DB}
}

func (r *galeriRepository) Save(ctx context.Context, galeri models.GaleriMitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, galeri)

	if err != nil {
		return result, err
	}

	return result, nil
}

// get all image