package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KTPMitraRepository interface {
	Save(ctx context.Context, eKTPModel models.KTPMitra) (*mongo.InsertOneResult, error)
	GetByMitraId(ctx context.Context, mitraId primitive.ObjectID) (models.KTPMitra, error)
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

func (r *ktpMitraRepository) GetByMitraId(ctx context.Context, mitraId primitive.ObjectID) (models.KTPMitra, error){
	var ktp models.KTPMitra
	err := r.DB.FindOne(ctx, bson.M{"mitra_id" : mitraId}).Decode(&ktp)

	if err != nil {
		return ktp, err
	}

	return ktp, nil
}

