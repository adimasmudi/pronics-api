package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type KomentarRepository interface {
	Save(ctx context.Context, newKomentar models.Komentar) (*mongo.InsertOneResult, error)
	GetAllByMitraId(ctx context.Context, mitraId primitive.ObjectID) ([]models.Komentar, error)
	GetById(ctx context.Context, komentarId primitive.ObjectID) (models.Komentar, error)
	GetByOrderId(ctx context.Context, orderId primitive.ObjectID) (models.Komentar, error)
	Update(ctx context.Context, komentarId primitive.ObjectID, newKomentar primitive.M) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, komentarId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type komentarRepository struct{
	DB *mongo.Collection
}

func NewKomentarRepository(DB *mongo.Collection) *komentarRepository{
	return &komentarRepository{DB}
}

func (r *komentarRepository) Save(ctx context.Context, newKomentar models.Komentar) (*mongo.InsertOneResult, error){
	result,err := r.DB.InsertOne(ctx, newKomentar)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *komentarRepository) GetAllByMitraId(ctx context.Context, mitraId primitive.ObjectID) ([]models.Komentar, error){
	var comments []models.Komentar

	currentRes, err := r.DB.Find(ctx, bson.M{"mitra_id" : mitraId})

	if err != nil{
		return comments, err
	}

	for currentRes.Next(ctx) {
        // looping to get each data and append to array
        var comment models.Komentar
        err := currentRes.Decode(&comment)
        if err != nil {
            return comments, err
        }

        comments =append(comments, comment)
    }

	if err := currentRes.Err(); err != nil {
        return comments, err
    }

	currentRes.Close(ctx)

	return comments, nil
}

func (r *komentarRepository) GetById(ctx context.Context, komentarId primitive.ObjectID) (models.Komentar, error){
	var komentar models.Komentar
	err := r.DB.FindOne(ctx, bson.M{"_id": komentarId}).Decode(&komentar)

	if err != nil{
		return komentar, err
	}

	return komentar, nil
}

func (r *komentarRepository) GetByOrderId(ctx context.Context, orderId primitive.ObjectID) (models.Komentar, error){
	var komentar models.Komentar
	err := r.DB.FindOne(ctx, bson.M{"order_id": orderId}).Decode(&komentar)

	if err != nil{
		return komentar, err
	}

	return komentar, nil
}

func (r *komentarRepository) Update(ctx context.Context, komentarId primitive.ObjectID, newKomentar primitive.M) (*mongo.UpdateResult, error){
	result, err := r.DB.UpdateOne(ctx,bson.M{"_id":komentarId},bson.M{"$set" : newKomentar})

	if err != nil{
		return result, err
	}

	return result, nil
}

func (r *komentarRepository) Delete(ctx context.Context, komentarId primitive.ObjectID) (*mongo.DeleteResult, error){
	result, err := r.DB.DeleteOne(ctx,bson.M{"_id":komentarId})

	if err != nil{
		return result, err
	}

	return result, nil
}