package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GaleriRepository interface {
	Save(ctx context.Context, galeri models.GaleriMitra) (*mongo.InsertOneResult, error)
	GetAllByIdMitra(ctx context.Context, mitraId primitive.ObjectID) ([]models.GaleriMitra, error)
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
func (r *galeriRepository) GetAllByIdMitra(ctx context.Context, mitraId primitive.ObjectID) ([]models.GaleriMitra, error){
	var galeriImages []models.GaleriMitra

	currentRes, err := r.DB.Find(ctx, bson.D{{"mitra_id", mitraId}})

	if err != nil{
		return galeriImages, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var galeriMitra models.GaleriMitra
        err := currentRes.Decode(&galeriMitra)
        if err != nil {
            return galeriImages, err
        }

        galeriImages =append(galeriImages, galeriMitra)
    }

	if err := currentRes.Err(); err != nil {
        return galeriImages, err
    }

	currentRes.Close(ctx)

	return galeriImages, nil
}