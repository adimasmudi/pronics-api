package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pronics-api/constants"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Login(ctx context.Context, input inputs.LoginUserInput) (models.User, string, error)
	Register(ctx context.Context, input inputs.RegisterUserInput) (*mongo.InsertOneResult, error)
	RegisterMitra(ctx context.Context, input inputs.RegisterMitraInput, fileName string) (*mongo.InsertOneResult, error)
	Signup(ctx context.Context, googleUser helper.GoogleUser) (string,error)
}

type userService struct {
	userRepository repositories.UserRepository
	customerRepository repositories.CustomerRepository
	mitraRepository repositories.MitraRepository
	rekeningRepository repositories.RekeningRepository
	ktpMitraRepository repositories.KTPMitraRepository
}

func NewUserService(userRepository repositories.UserRepository, customerRepository repositories.CustomerRepository, mitraRepository repositories.MitraRepository, rekeningRepository repositories.RekeningRepository, ktpMitraRepository repositories.KTPMitraRepository) *userService{
	return &userService{userRepository, customerRepository ,mitraRepository, rekeningRepository, ktpMitraRepository}
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

	if input.Type != constants.UserCustomer && input.Type != constants.UserMitra{
		return nil, errors.New("tipe user hanya boleh customer atau mitra")
	}


	newUser := models.User{
		ID : primitive.NewObjectID(),
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

	newCustomer := models.Customer{
		ID : primitive.NewObjectID(),
		UserId: registeredUser.InsertedID.(primitive.ObjectID),
		Username : input.Email,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	registeredCustomer, err := s.customerRepository.SaveRegisterUser(ctx, newCustomer)

	if err != nil{
		return registeredCustomer, err
	}

	return registeredUser, nil
}

func (s *userService) RegisterMitra(ctx context.Context, input inputs.RegisterMitraInput, fileName string) (*mongo.InsertOneResult, error){
	userExist, _ := s.userRepository.IsUserExist(ctx,input.Email )

	if userExist{
		return nil, errors.New("user with this email already exist")
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)

	if err != nil{
		return nil, err
	}

	if input.Type != constants.UserCustomer && input.Type != constants.UserMitra{
		return nil, errors.New("tipe user hanya boleh customer atau mitra")
	}

	if input.MitraType != constants.MitraIndividu && input.MitraType != constants.MitraToko{
		return nil, errors.New("tipe mitra hanya boleh antara individu atau toko")
	}

	newUser := models.User{
		ID : primitive.NewObjectID(),
		NamaLengkap: input.NamaLengkap,
		Email : input.Email,
		Password : string(passwordHash),
		NoTelepon: input.NomerTelepon,
		Deskripsi: input.Deskripsi,
		Type : input.Type,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	registeredUser, err := s.userRepository.Save(ctx,newUser)


	if err != nil{
		return nil, err
	}

	wilayahMitraId, err := primitive.ObjectIDFromHex(input.Wilayah)

	if err != nil{
		return nil, errors.New("id wilayah yang diterima")
	}

	bidangStrArr := input.Bidang
	bidangStr := strings.TrimSpace(bidangStrArr)
	bidangStr = strings.Replace(bidangStr, "[", "", -1)
	bidangStr = strings.Replace(bidangStr, "]", "", -1)
	bidangArr := strings.Split(bidangStr, ",")

	var bidangMitra []primitive.ObjectID

	for _, bidang := range bidangArr{
		eachBidang, _ := primitive.ObjectIDFromHex(bidang)
		bidangMitra = append(bidangMitra, eachBidang)
	}

	newMitra := models.Mitra{
		ID : primitive.NewObjectID(),
		UserId: registeredUser.InsertedID.(primitive.ObjectID),
		NamaToko: input.NamaToko,
		Alamat : input.Alamat,
		MitraType: input.MitraType,
		Status : constants.MitraInActive,
		Wilayah : wilayahMitraId,
		Bidang : bidangMitra,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	registeredMitra, err := s.mitraRepository.SaveMitra(ctx, newMitra)

	if err != nil{
		return registeredMitra, err
	}

	newRekening := models.Rekening{
		ID : primitive.NewObjectID(),
		UserId: registeredUser.InsertedID.(primitive.ObjectID),
		BankId : primitive.NewObjectID(), // generate id sementara
		NamaPemilik: input.NamaPemilikRekening,
		NomerRekening: input.NomerRekening,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	rekeningAdded, err := s.rekeningRepository.SaveRekening(ctx, newRekening)

	if err != nil{
		return rekeningAdded, err
	}

	newKTPMitra := models.KTPMitra{
		ID : primitive.NewObjectID(),
		MitraId: registeredMitra.InsertedID.(primitive.ObjectID),
		GambarKtp: os.Getenv("CLOUD_STORAGE_READ_LINK")+"ektp/"+fileName,
		CreatedAt: time.Now(),
		UpdatedAt : time.Now(),
	}

	ktpAdded, err := s.ktpMitraRepository.Save(ctx, newKTPMitra)

	if err != nil{
		return ktpAdded, err
	}

	return registeredUser, nil
}

func (s *userService) Login(ctx context.Context, input inputs.LoginUserInput) (models.User, string, error){

	user, err := s.userRepository.FindByEmail(ctx,input.Email)

	if err != nil{
		return user, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))

	if err != nil{
		return user, "", errors.New("wrong Password")
	}

	token, err := helper.GenerateToken(user.ID)

	if err != nil{
		return user, "", errors.New("can't generate token")
	}

	return user, token, nil
}

func (s *userService) Signup(ctx context.Context, googleUser helper.GoogleUser) (string,error){
	userExist, _ := s.userRepository.IsUserExist(ctx,googleUser.Email)

	if !userExist{
		newUser := models.User{
			ID : primitive.NewObjectID(),
			NamaLengkap: strings.Split(googleUser.Email, "@")[0],
			Email : googleUser.Email,
			Type : constants.UserCustomer,
			CreatedAt: time.Now(),
			UpdatedAt : time.Now(),
		}
	
		registeredUser, err := s.userRepository.Save(ctx,newUser)
	
	
		if err != nil{
			return "", err
		}
	
		newCustomer := models.Customer{
			ID : primitive.NewObjectID(),
			UserId: registeredUser.InsertedID.(primitive.ObjectID),
			Username : googleUser.Email,
			GambarCustomer : googleUser.Picture,
			CreatedAt: time.Now(),
			UpdatedAt : time.Now(),
		}
	
		registeredCustomer, err := s.customerRepository.SaveRegisterUser(ctx, newCustomer)
		
		fmt.Println(registeredCustomer)

		if err != nil{
			return "", err
		}

		token,err := helper.GenerateToken(registeredUser.InsertedID.(primitive.ObjectID))

		if err != nil{
			return token,err
		}

		return token, nil
	}

	userFound, err := s.userRepository.FindByEmail(ctx,googleUser.Email)

	
	if err != nil{
		return "", err
	}

	token,err := helper.GenerateToken(userFound.ID)

	if err != nil{
		return token, err
	}

	return token, nil
}
