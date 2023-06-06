package repositories


import (
	"context"
	"pronics-api/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type MitraRepository interface {
	FindByEmail(ctx context.Context, email string) (models.User, error)
	Save(ctx context.Context, mitra models.Mitra) (*mongo.InsertOneResult, error)
	IsUserExist(ctx context.Context, email string) (bool, error)
}

type mitraRepository struct{
	DB *mongo.Collection
}

func NewMitraRepository(DB *mongo.Collection) *mitraRepository{
	return &mitraRepository{DB}
}