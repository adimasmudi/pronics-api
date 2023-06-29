package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderDetailRepository interface {
	Save(ctx context.Context, orderDetail models.OrderDetail) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, IdOrderDetail primitive.ObjectID, newOrderDetail primitive.M) (*mongo.UpdateResult, error)
	GetById(ctx context.Context, Id primitive.ObjectID) (models.OrderDetail,  error)
	GetByOrderId(ctx context.Context, orderId primitive.ObjectID) (models.OrderDetail,  error)
}

type orderDetailRepository struct{
	DB *mongo.Collection
}

func NewOrderDetailRepository(DB *mongo.Collection) *orderDetailRepository{
	return &orderDetailRepository{DB}
}

func (r *orderDetailRepository) Save(ctx context.Context, orderDetail models.OrderDetail) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, orderDetail)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *orderDetailRepository) GetById(ctx context.Context, Id primitive.ObjectID) (models.OrderDetail,  error){

	var orderDetail models.OrderDetail

	err := r.DB.FindOne(ctx, bson.M{"_id": Id}).Decode(&orderDetail)

	if err != nil{
		return orderDetail, err
	}

	return orderDetail, nil
}

func (r *orderDetailRepository) GetByOrderId(ctx context.Context, orderId primitive.ObjectID) (models.OrderDetail,  error){

	var orderDetail models.OrderDetail

	err := r.DB.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&orderDetail)

	if err != nil{
		return orderDetail, err
	}

	return orderDetail, nil
}

func (r *orderDetailRepository) Update(ctx context.Context, IdOrderDetail primitive.ObjectID, newOrderDetail primitive.M) (*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx, bson.M{"_id" : IdOrderDetail},bson.M{"$set" : newOrderDetail})

	if err != nil{
		return result, err
	}

	return result, nil
}