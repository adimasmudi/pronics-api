package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RekeningRepository interface {
	SaveRekening(ctx context.Context, rekening models.Rekening) (*mongo.InsertOneResult, error)
	GetRekeningByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Rekening,  error)
}

type rekeningRepository struct{
	DB *mongo.Collection
}

func NewRekeningRepository(DB *mongo.Collection) *rekeningRepository{
	return &rekeningRepository{DB}
}

func (r *rekeningRepository) SaveRekening(ctx context.Context, rekening models.Rekening) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, rekening)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *rekeningRepository) GetRekeningByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Rekening,  error){

	var rekening models.Rekening

	err := r.DB.FindOne(ctx, bson.M{"user_id": IdUser}).Decode(&rekening)

	if err != nil{
		return rekening, err
	}

	return rekening, nil
}