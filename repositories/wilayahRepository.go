package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type WilayahRepository interface {
	SaveWilayah(ctx context.Context, wilayah models.WilayahCakupan) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.WilayahCakupan, error)
	FindById(ctx context.Context, ID primitive.ObjectID) (models.WilayahCakupan, error)
}

type wilayahRepository struct{
	DB *mongo.Collection
}

func NewWilayahRepository(DB *mongo.Collection) *wilayahRepository{
	return &wilayahRepository{DB}
}

func (r *wilayahRepository) SaveWilayah(ctx context.Context, wilayah models.WilayahCakupan) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, wilayah)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *wilayahRepository) FindAll(ctx context.Context) ([]models.WilayahCakupan, error){
	var wilayahCakupans []models.WilayahCakupan

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return wilayahCakupans, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var wilayahCakupan models.WilayahCakupan
        err := currentRes.Decode(&wilayahCakupan)
        if err != nil {
            return wilayahCakupans, err
        }

        wilayahCakupans =append(wilayahCakupans, wilayahCakupan)
    }

	if err := currentRes.Err(); err != nil {
        return wilayahCakupans, err
    }

	currentRes.Close(ctx)

	return wilayahCakupans, nil
}

func (r *wilayahRepository) FindById(ctx context.Context, ID primitive.ObjectID) (models.WilayahCakupan, error){
	var wilayah models.WilayahCakupan

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&wilayah)
	
	if err != nil{
		return wilayah, err
	}

	return wilayah, nil
}