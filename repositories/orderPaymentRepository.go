package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderPaymentRepository interface {
	Save(ctx context.Context, orderPayment models.OrderPayment) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, IdOrderPayment primitive.ObjectID, newOrderPayment primitive.M) (*mongo.UpdateResult, error)
	GetById(ctx context.Context, Id primitive.ObjectID) (models.OrderPayment,  error)
	GetByOrderDetailId(ctx context.Context, orderId primitive.ObjectID) (models.OrderPayment,  error)
}

type orderPaymentRepository struct{
	DB *mongo.Collection
}

func NewOrderPaymentRepository(DB *mongo.Collection) *orderPaymentRepository{
	return &orderPaymentRepository{DB}
}

func (r *orderPaymentRepository) Save(ctx context.Context, orderPayment models.OrderPayment) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, orderPayment)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *orderPaymentRepository) GetById(ctx context.Context, Id primitive.ObjectID) (models.OrderPayment,  error){

	var orderPayment models.OrderPayment

	err := r.DB.FindOne(ctx, bson.M{"_id": Id}).Decode(&orderPayment)

	if err != nil{
		return orderPayment, err
	}

	return orderPayment, nil
}

func (r *orderPaymentRepository) GetByOrderDetailId(ctx context.Context, orderDetailId primitive.ObjectID) (models.OrderPayment,  error){

	var orderPayment models.OrderPayment

	err := r.DB.FindOne(ctx, bson.M{"order_detail_id": orderDetailId}).Decode(&orderPayment)

	if err != nil{
		return orderPayment, err
	}

	return orderPayment, nil
}

func (r *orderPaymentRepository) Update(ctx context.Context, IdOrderPayment primitive.ObjectID, newOrderPayment primitive.M) (*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx, bson.M{"_id" : IdOrderPayment},bson.M{"$set" : newOrderPayment})

	if err != nil{
		return result, err
	}

	return result, nil
}