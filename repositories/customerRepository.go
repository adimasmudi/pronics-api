package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type CustomerRepository interface {
	SaveRegisterUser(ctx context.Context, customer models.Customer) (*mongo.InsertOneResult, error)
	GetCustomerById(ctx context.Context, ID primitive.ObjectID) (models.Customer,  error)
	GetCustomerByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Customer,  error)
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

func (r *customerRepository) GetCustomerById(ctx context.Context, ID primitive.ObjectID) (models.Customer,  error){

	var customer models.Customer

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&customer)

	if err != nil{
		return customer, err
	}

	return customer, nil
}

func (r *customerRepository) GetCustomerByIdUser(ctx context.Context, IdUser primitive.ObjectID) (models.Customer,  error){

	var customer models.Customer

	err := r.DB.FindOne(ctx, bson.M{"user_id": IdUser}).Decode(&customer)

	if err != nil{
		return customer, err
	}

	return customer, nil
}