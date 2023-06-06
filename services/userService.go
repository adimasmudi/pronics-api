package services

import (
	"context"
	"errors"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, input inputs.LoginUserInput) (models.User, string, error)
	Register(ctx context.Context, input inputs.RegisterUserInput) (*mongo.InsertOneResult, error)
}

type userService struct {
	userRepository repositories.UserRepository
}

func NewUserService(userRepository repositories.UserRepository) *userService{
	return &userService{userRepository}
}

func (s *userService) Register(ctx context.Context, input inputs.RegisterUserInput) (*mongo.InsertOneResult, error){
	
	userExist, _ := s.userRepository.IsUserExist(ctx,input.Email )

	if userExist{
		return nil, errors.New("user with this email already exist")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil{
		return nil, err
	}


	newUser := models.User{
		NamaLengkap: input.NamaLengkap,
		Email : input.Email,
		Password : string(passwordHash),
		Type : input.Type,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	registeredUser, err := s.userRepository.Save(ctx,newUser)

	if err != nil{
		return nil, err
	}

	return registeredUser, nil
}

func (s *userService) Login(ctx context.Context, input inputs.LoginUserInput) (models.User, string, error){

	user, err := s.userRepository.FindByEmail(ctx,input.Email)

	if err != nil{
		return user, "", errors.New("email not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil{
		return user, "", errors.New("wrong Password")
	}

	token, err := helper.GenerateToken(user.Email)

	if err != nil{
		return user, "", errors.New("can't generate token")
	}

	return user, token, nil
}
