package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type LayananMitraRepository interface {
	Save(ctc context.Context, layananMitra models.LayananMitra) (*mongo.InsertOneResult, error)
	FindAllByBidangAndMitra(ctx context.Context, bidangId primitive.ObjectID, mitraId primitive.ObjectID) ([]models.LayananMitra, error)
	GetById(ctx context.Context, ID primitive.ObjectID) (models.LayananMitra, error)
	DeleteLayananMitra(ctx context.Context, IdLayananMitra primitive.ObjectID) (*mongo.DeleteResult, error)
}

type layananMitraRepository struct{
	DB *mongo.Collection
}

func NewLayananMitraRepository(DB *mongo.Collection) *layananMitraRepository{
	return &layananMitraRepository{DB}
}

// save layanan mitra
func (r *layananMitraRepository) Save(ctx context.Context, layananMitra models.LayananMitra) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, layananMitra)

	if err != nil {
		return result, err
	}

	return result, nil
}

// find all
func (r *layananMitraRepository) FindAllByBidangAndMitra(ctx context.Context, bidangId primitive.ObjectID, mitraId primitive.ObjectID) ([]models.LayananMitra, error){
	var layananMitras []models.LayananMitra

	currentRes, err := r.DB.Find(ctx, bson.M{"bidang_id": bidangId,"mitra_id": mitraId})

	if err != nil{
		return layananMitras, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var layananMitra models.LayananMitra
        err := currentRes.Decode(&layananMitra)
        if err != nil {
            return layananMitras, err
        }

        layananMitras =append(layananMitras, layananMitra)
    }

	if err := currentRes.Err(); err != nil {
        return layananMitras, err
    }

	currentRes.Close(ctx)

	return layananMitras, nil
}

// get by id
func (r *layananMitraRepository) GetById(ctx context.Context, ID primitive.ObjectID) (models.LayananMitra, error){
	var layananMitra models.LayananMitra
	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&layananMitra)

	if err != nil{
		return layananMitra, err
	}

	return layananMitra, nil
}

// delete
func (r *layananMitraRepository) DeleteLayananMitra(ctx context.Context, IdLayananMitra primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":IdLayananMitra})

	if err != nil{
		return result, err
	}

	return result, nil
}

