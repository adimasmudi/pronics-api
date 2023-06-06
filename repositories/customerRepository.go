package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository interface {
	FindByEmail(ctx context.Context, email string) (models.Customer, error)
	Save(ctx context.Context, customer models.Customer) (*mongo.InsertOneResult, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type customerRepository struct{
	DB *mongo.Collection
}

func NewcustomerRepository(DB *mongo.Collection) *customerRepository{
	return &customerRepository{DB}
}