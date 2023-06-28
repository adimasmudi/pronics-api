package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	Create(ctx context.Context, newOrder models.Order) (*mongo.InsertOneResult, error)
	GetById(ctx context.Context, ID primitive.ObjectID) (models.Order, error)
	GetOrderByCustomerIdNMitraId(ctx context.Context, customerId primitive.ObjectID, mitraId primitive.ObjectID) (models.Order, error)
	GetAllOrder(ctx context.Context, customerId primitive.ObjectID)([]models.Order, error)
	GetAllOrderCustomer(ctx context.Context, customerId primitive.ObjectID) ([]models.Order, error)
	GetAllOrderMitra(ctx context.Context, mitraId primitive.ObjectID) ([]models.Order, error)
	UpdateOrder(ctx context.Context, ID primitive.ObjectID, newOrder primitive.M)(*mongo.UpdateResult, error)
}

type orderRepository struct{
	DB *mongo.Collection
}

func NewOrderRepository(DB *mongo.Collection) *orderRepository{
	return &orderRepository{DB}
}

func (r *orderRepository) Create(ctx context.Context, newOrder models.Order) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, newOrder)

	if err != nil {
		return result, err
	}

	return result, nil
}

// get order by id
func (r *orderRepository) GetById(ctx context.Context, ID primitive.ObjectID) (models.Order, error){
	var order models.Order

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&order)

	if err != nil{
		return order, err
	}

	return order, nil
}

// get order based on id customer and id mitra
func (r *orderRepository) GetOrderByCustomerIdNMitraId(ctx context.Context, customerId primitive.ObjectID, mitraId primitive.ObjectID) (models.Order, error){
	var order models.Order

	err := r.DB.FindOne(ctx, bson.M{"customer_id": customerId,"mitra_id":mitraId}).Decode(&order)

	if err != nil{
		return order, err
	}

	return order, nil
}

// get all order
func (r *orderRepository) GetAllOrder(ctx context.Context) ([]models.Order, error){
	var orders []models.Order
	
	currentRes, err := r.DB.Find(ctx, bson.D{})

	if err != nil{
		return nil, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var order models.Order
        err := currentRes.Decode(&order)
        if err != nil {
            return orders, err
        }

        orders =append(orders, order)
    }

	if err := currentRes.Err(); err != nil {
        return orders, err
    }

	currentRes.Close(ctx)

	return orders, nil
}

// get all order based on id customer
func (r *orderRepository) GetAllOrderCustomer(ctx context.Context, customerId primitive.ObjectID) ([]models.Order, error){
	var orders []models.Order
	
	currentRes, err := r.DB.Find(ctx, bson.M{"customer_id":customerId})

	if err != nil{
		return nil, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var order models.Order
        err := currentRes.Decode(&order)
        if err != nil {
            return orders, err
        }

        orders =append(orders, order)
    }

	if err := currentRes.Err(); err != nil {
        return orders, err
    }

	currentRes.Close(ctx)

	return orders, nil
}

// get all order based on id mitra
func (r *orderRepository) GetAllOrderMitra(ctx context.Context, mitraId primitive.ObjectID) ([]models.Order, error){
	var orders []models.Order
	
	currentRes, err := r.DB.Find(ctx, bson.M{"mitra_id":mitraId})

	if err != nil{
		return nil, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var order models.Order
        err := currentRes.Decode(&order)
        if err != nil {
            return orders, err
        }

        orders =append(orders, order)
    }

	if err := currentRes.Err(); err != nil {
        return orders, err
    }

	currentRes.Close(ctx)

	return orders, nil
}

// update order
func (r *orderRepository) UpdateOrder(ctx context.Context, ID primitive.ObjectID, newOrder primitive.M)(*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":ID},bson.M{"$set" : newOrder})

	if err != nil{
		return result, err
	}

	return result, nil
}
