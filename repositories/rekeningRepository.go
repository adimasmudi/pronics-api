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
	UpdateRekening(ctx context.Context, ID primitive.ObjectID, newRekening primitive.M) (*mongo.UpdateResult, error)
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

func (r *rekeningRepository) UpdateRekening(ctx context.Context, ID primitive.ObjectID, newRekening primitive.M) (*mongo.UpdateResult, error){
	data, err := r.DB.UpdateOne(
		ctx,
		bson.M{"_id" : ID},
		bson.M{"$set" : newRekening},
	)

	if err != nil{
		return data, err
	}

	return data, nil
}