package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LayananRepository interface {
	SaveLayanan(ctx context.Context, layanan models.Layanan) (*mongo.InsertOneResult, error)
	FindAll(ctx context.Context) ([]models.Layanan, error)
	GetById(ctx context.Context, ID primitive.ObjectID) (models.Layanan, error)
	FindAllByBidangId(ctx context.Context, bidangId primitive.ObjectID) ([]models.Layanan, error)
	DeleteLayanan(ctx context.Context, IdLayanan primitive.ObjectID) (*mongo.DeleteResult, error)
}

type layananRepository struct{
	DB *mongo.Collection
}

func NewLayananRepository(DB *mongo.Collection) *layananRepository{
	return &layananRepository{DB}
}

// save layanan
func (r *layananRepository) SaveLayanan(ctx context.Context, layanan models.Layanan) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, layanan)

	if err != nil {
		return result, err
	}

	return result, nil
}

// find all
func (r *layananRepository) FindAll(ctx context.Context) ([]models.Layanan, error){
	var layanans []models.Layanan

	currentRes, err := r.DB.Find(ctx, bson.D{{}})

	if err != nil{
		return layanans, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var layanan models.Layanan
        err := currentRes.Decode(&layanan)
        if err != nil {
            return layanans, err
        }

        layanans =append(layanans, layanan)
    }

	if err := currentRes.Err(); err != nil {
        return layanans, err
    }

	currentRes.Close(ctx)

	return layanans, nil
}

// get by id
func (r *layananRepository) GetById(ctx context.Context, ID primitive.ObjectID) (models.Layanan, error){
	var layanan models.Layanan
	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&layanan)

	if err != nil{
		return layanan, err
	}

	return layanan, nil
}

// get all by bidang
func (r *layananRepository) FindAllByBidangId(ctx context.Context, bidangId primitive.ObjectID) ([]models.Layanan, error){
	var layanans []models.Layanan

	currentRes, err := r.DB.Find(ctx, bson.D{{"bidang_id",bidangId}})

	if err != nil{
		return layanans, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var layanan models.Layanan
        err := currentRes.Decode(&layanan)
        if err != nil {
            return layanans, err
        }

        layanans =append(layanans, layanan)
    }

	if err := currentRes.Err(); err != nil {
        return layanans, err
    }

	currentRes.Close(ctx)

	return layanans, nil
}

// delete
func (r *layananRepository) DeleteLayanan(ctx context.Context, IdLayanan primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":IdLayanan})

	if err != nil{
		return result, err
	}

	return result, nil
}
