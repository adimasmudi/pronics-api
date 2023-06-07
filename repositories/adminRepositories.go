package repositories

import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AdminRepository interface{
	FindByEmail(ctx context.Context,email string) (models.Admin, error)
	Save(ctx context.Context,admin models.Admin) (*mongo.InsertOneResult, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
	GetAdminById(ctx context.Context, ID primitive.ObjectID) (models.Admin, error)
}

type adminRepository struct{
	DB *mongo.Collection
}

func NewAdminRepository(DB *mongo.Collection) *adminRepository{
	return &adminRepository{DB}
}

func (r *adminRepository) Save(ctx context.Context,admin models.Admin) (*mongo.InsertOneResult, error) {
	r.DB.Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys : bson.D{{Key: "email", Value: 1},{Key:"username", Value:1}},
			Options : options.Index().SetUnique(true),
		},
	)
	
	result,err := r.DB.InsertOne(ctx, admin)

	if err != nil {
		return result, err
	}

	return result, nil
}

func (r *adminRepository) FindByEmail(ctx context.Context, email string) (models.Admin,  error){

	var admin models.Admin

	err := r.DB.FindOne(ctx, bson.M{"email": email}).Decode(&admin)

	if err != nil{
		return admin, err
	}

	return admin, nil
}

func (r *adminRepository) IsUserExist(ctx context.Context, email string) (bool, error){
	var admin models.Admin

	err := r.DB.FindOne(ctx, bson.M{"email": email}).Decode(&admin)

	if err != nil{
		return false, err
	}

	return true, nil
}

func (r *adminRepository) GetAdminById(ctx context.Context, ID primitive.ObjectID) (models.Admin, error){
	var admin models.Admin

	err := r.DB.FindOne(ctx, bson.M{"_id": ID}).Decode(&admin)

	if err != nil{
		return admin, err
	}

	return admin, nil
}