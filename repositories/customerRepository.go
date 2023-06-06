package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository interface {
	SaveRegisterUser(ctx context.Context, customer models.Customer) (*mongo.InsertOneResult, error)
}

type customerRepository struct{
	DB *mongo.Collection
}

func NewcustomerRepository(DB *mongo.Collection) *customerRepository{
	return &customerRepository{DB}
}

func (r *customerRepository) SaveRegisterUser(ctx context.Context, customer models.Customer) (*mongo.InsertOneResult, error){
	
	result,err := r.DB.InsertOne(ctx, customer)

	if err != nil {
		return result, err
	}

	return result, nil
}