package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BankRepository interface {
	Save(ctx context.Context, bank models.Bank) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Bank, error)
	GetBankById(ctx context.Context, Id primitive.ObjectID) (models.Bank,  error)
}

type bankRepository struct{
	DB *mongo.Collection
}

func NewBankRepository(DB *mongo.Collection) *bankRepository{
	return &bankRepository{DB}
}

func (r *bankRepository) Save(ctx context.Context, bank models.Bank) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, bank)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *bankRepository) FindAll(ctx context.Context) ([]models.Bank, error){
	var banks []models.Bank

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return banks, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var bank models.Bank
        err := currentRes.Decode(&bank)
        if err != nil {
            return banks, err
        }

        banks =append(banks, bank)
    }

	if err := currentRes.Err(); err != nil {
        return banks, err
    }

	currentRes.Close(ctx)

	return banks, nil
}

func (r *bankRepository) GetBankById(ctx context.Context, Id primitive.ObjectID) (models.Bank,  error){

	var bank models.Bank

	err := r.DB.FindOne(ctx, bson.M{"_id": Id}).Decode(&bank)

	if err != nil{
		return bank, err
	}

	return bank, nil
}