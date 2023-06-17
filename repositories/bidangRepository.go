package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BidangRepository interface {
	SaveBidang(ctx context.Context, bidang models.Bidang) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Bidang, error)
	GetById(ctx context.Context, ID primitive.ObjectID) (models.Bidang, error)
	UpdateBidang(ctx context.Context, IdBidang primitive.ObjectID, newBidang primitive.M) (*mongo.UpdateResult, error)
}

type bidangRepository struct{
	DB *mongo.Collection
}

func NewBidangRepository(DB *mongo.Collection) *bidangRepository{
	return &bidangRepository{DB}
}

func (r *bidangRepository) SaveBidang(ctx context.Context, bidang models.Bidang) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, bidang)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *bidangRepository) FindAll(ctx context.Context) ([]models.Bidang, error){
	var bidangs []models.Bidang

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return bidangs, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var bidang models.Bidang
        err := currentRes.Decode(&bidang)
        if err != nil {
            return bidangs, err
        }

        bidangs =append(bidangs, bidang)
    }

	if err := currentRes.Err(); err != nil {
        return bidangs, err
    }

	currentRes.Close(ctx)

	return bidangs, nil
}

func (r *bidangRepository) GetById(ctx context.Context, ID primitive.ObjectID) (models.Bidang, error){
	var bidang models.Bidang
	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&bidang)

	if err != nil{
		return bidang, err
	}

	return bidang, nil
}

func (r *bidangRepository) UpdateBidang(ctx context.Context, IdBidang primitive.ObjectID, newBidang primitive.M) (*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":IdBidang},bson.M{"$set" : newBidang})

	if err != nil{
		return result, err
	}

	return result, nil
}

func (r *bidangRepository) DeleteBidang(ctx context.Context, IdBidang primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":IdBidang})

	if err != nil{
		return result, err
	}

	return result, nil
}