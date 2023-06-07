package services

import (
	"context"
	"errors"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AdminService interface {
	Register(ctx context.Context, input inputs.RegisterAdminInput) (*mongo.InsertOneResult, error)
	Login(ctx context.Context, input inputs.LoginAdminInput) (models.Admin, string, error)
	GetAdminProfile(ctx context.Context, ID primitive.ObjectID) (models.Admin, error)
}

type adminService struct{
	adminRepository repositories.AdminRepository
}

func NewAdminService(adminRepository repositories.AdminRepository) *adminService{
	return &adminService{adminRepository}
}

func (s *adminService) Register(ctx context.Context, input inputs.RegisterAdminInput) (*mongo.InsertOneResult, error){
	
	userExist, _ := s.adminRepository.IsUserExist(ctx,input.Email )

	if userExist{
		return nil, errors.New("admin with this email already exist")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil{
		return nil, err
	}


	newAdmin := models.Admin{
		ID : primitive.NewObjectID(),
		Username : input.Username,
		Email : input.Email,
		Password : string(passwordHash),
		IsAdmin: true,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	registeredAdmin, err := s.adminRepository.Save(ctx,newAdmin)

	if err != nil{
		return nil, err
	}

	return registeredAdmin, nil
}

func (s *adminService) Login(ctx context.Context, input inputs.LoginAdminInput) (models.Admin, string, error){

	admin, err := s.adminRepository.FindByEmail(ctx,input.Email)

	if err != nil{
		return admin, "", errors.New("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(input.Password))

	if err != nil{
		return admin, "", errors.New("wrong Password")
	}

	token, err := helper.GenerateToken(admin.ID)

	if err != nil{
		return admin, "", errors.New("can't generate token")
	}

	return admin, token, nil
}

func (s *adminService) GetAdminProfile(ctx context.Context, ID primitive.ObjectID) (models.Admin, error){
	admin, err := s.adminRepository.GetAdminById(ctx, ID)

	if err != nil{
		return admin, err
	}

	return admin, nil
}