package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MitraRepository interface {
	SaveMitra(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error)
	GetMitraById(ctx context.Context, ID primitive.ObjectID) (models.Mitra,  error)
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

func (r *mitraRepository) GetMitraById(ctx context.Context, ID primitive.ObjectID) (models.Mitra,  error){

	var mitra models.Mitra

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&mitra)

	if err != nil{
		return mitra, err
	}

	return mitra, nil
}