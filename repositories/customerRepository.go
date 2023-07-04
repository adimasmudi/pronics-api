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
	UpdateProfil(ctx context.Context, ID primitive.ObjectID, newCustomer primitive.M)(*mongo.UpdateResult, error)
	UpdateAlamatCustomer(ctx context.Context, customerId primitive.ObjectID, newAlamatInCustomer primitive.M) (*mongo.UpdateResult, error)
	GetAllCustomer(ctx context.Context) ([]models.Customer, error)
}

type customerRepository struct{
	DB *mongo.Collection
}

func NewCustomerRepository(DB *mongo.Collection) *customerRepository{
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

func (r *customerRepository) UpdateProfil(ctx context.Context, ID primitive.ObjectID, newCustomer primitive.M)(*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":ID},bson.M{"$set" : newCustomer})

	if err != nil{
		return result, err
	}

	return result, nil
}

func (r *customerRepository) UpdateAlamatCustomer(ctx context.Context, customerId primitive.ObjectID, newAlamatInCustomer primitive.M) (*mongo.UpdateResult, error){
	data, err := r.DB.UpdateOne(
		ctx,
		bson.M{"_id" : customerId},
		bson.M{"$set" : newAlamatInCustomer},
	)

	if err != nil{
		return data, err
	}

	return data, nil
}

func (r *customerRepository) GetAllCustomer(ctx context.Context) ([]models.Customer, error){
	var customers []models.Customer

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return customers, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var customer models.Customer
        err := currentRes.Decode(&customer)
        if err != nil {
            return customers, err
        }

        customers =append(customers, customer)
    }

	if err := currentRes.Err(); err != nil {
        return customers, err
    }

	currentRes.Close(ctx)

	return customers, nil
}