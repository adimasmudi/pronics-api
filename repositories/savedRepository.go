package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SavedRepository interface {
	Save(ctx context.Context, newSaved models.Saved) (*mongo.InsertOneResult, error)
	GetAll(ctx context.Context, customerId primitive.ObjectID) ([]models.Saved, error)
	Delete(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error)
	GetByIdMitra(ctx context.Context, mitraId primitive.ObjectID) (models.Saved, error)
	GetByIdCustomerNMitra(ctx context.Context, customerId primitive.ObjectID, mitraId primitive.ObjectID) (models.Saved, error)
}

type savedRepository struct{
	DB *mongo.Collection
}

func NewSavedRepository(DB *mongo.Collection) *savedRepository{
	return &savedRepository{DB}
}

func (r *savedRepository) Save(ctx context.Context, newSaved models.Saved) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, newSaved)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *savedRepository) GetAll(ctx context.Context, customerId primitive.ObjectID) ([]models.Saved, error){
	var allSaveds []models.Saved
	
	currentRes, err := r.DB.Find(ctx,bson.M{"customer_id" : customerId})

	if err != nil{
		return nil, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var Saved models.Saved
        err := currentRes.Decode(&Saved)
        if err != nil {
            return allSaveds, err
        }

        allSaveds =append(allSaveds, Saved)
    }

	if err := currentRes.Err(); err != nil {
        return allSaveds, err
    }

	currentRes.Close(ctx)

	return allSaveds, nil
}

func (r *savedRepository) Delete(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":savedId})

	if err != nil{
		return result, err
	}

	return result, nil
}

func (r *savedRepository) GetByIdMitra(ctx context.Context, mitraId primitive.ObjectID) (models.Saved, error){
	var saved models.Saved
	err := r.DB.FindOne(ctx, bson.M{"mitra_id": mitraId}).Decode(&saved)

	if err != nil{
		return saved, err
	}

	return saved, nil
}

func (r *savedRepository) GetByIdCustomerNMitra(ctx context.Context, customerId primitive.ObjectID, mitraId primitive.ObjectID) (models.Saved, error){
	var saved models.Saved
	err := r.DB.FindOne(ctx, bson.M{"customer_id" : customerId,"mitra_id": mitraId}).Decode(&saved)

	if err != nil{
		return saved, err
	}

	return saved, nil
}