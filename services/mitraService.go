package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MitraService interface {
	GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error)
	UpdateProfileMitra(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilMitraInput, fileName string) (*mongo.UpdateResult, error)
	UploadGaleriImage(ctx context.Context, ID primitive.ObjectID, fileNames []string) (*mongo.UpdateResult, error)
}

type mitraService struct {
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
	galeriMitraRepository repositories.GaleriRepository
}

func NewMitraService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, galeriMitraRepository repositories.GaleriRepository) *mitraService{
	return &mitraService{userRepository, mitraRepository, galeriMitraRepository}
}

func (s *mitraService) GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error){
	var data formatters.MitraResponse

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, user.ID)

	if err != nil{
		return data, err
	}

	data = helper.MapperMitra(user, mitra)

	return data, nil
}

func (s *mitraService) UpdateProfileMitra(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilMitraInput, fileName string) (*mongo.UpdateResult, error){
	var newMitra primitive.M
	
	if fileName != ""{
		newMitra = bson.M{
			"namatoko" : input.NamaToko,
			"gambarmitra": os.Getenv("CLOUD_STORAGE_READ_LINK")+"mitra/"+fileName,
			"alamat" : input.Alamat,
			"updatedat" : time.Now(),
		}
	}else{
		newMitra = bson.M{
			"namatoko" : input.NamaToko,
			"alamat" : input.Alamat,
			"updatedat" : time.Now(),
		}
	}
	

	newUser := bson.M{
		"namalengkap" : input.NamaLengkap,
		"email" : input.Email,
		"notelepon" : input.NoHandphone,
		"deskripsi" : input.Deskripsi,
		"jeniskelamin" : input.JenisKelamin,
		"tanggallahir" : input.TanggalLahir,
		"updatedat": time.Now(),
	}

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx,ID)

	if err != nil{
		return nil, err
	}

	updatedUser, err := s.userRepository.UpdateUser(ctx, ID, newUser)

	if err != nil{
		return nil, err
	}

	updatedMitra, err := s.mitraRepository.UpdateProfil(ctx, mitra.ID,newMitra)

	if err != nil{
		return nil, err
	}

	fmt.Println(updatedMitra)

	return updatedUser, nil
}

func (s *mitraService) UploadGaleriImage(ctx context.Context, ID primitive.ObjectID, fileNames []string)(*mongo.UpdateResult, error){
	var newGaleriMitras []primitive.ObjectID

	if len(fileNames) == 0{
		return nil, errors.New("Tidak ada gambar di upload")
	}

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx,ID)

	if len(mitra.GaleriMitra) > 0{
		newGaleriMitras = append(newGaleriMitras, mitra.GaleriMitra...)
	}

	if err != nil{
		return nil, err
	}

	for _, fileName := range(fileNames){
		newGaleriMitra := models.GaleriMitra{
			ID : primitive.NewObjectID(),
			MitraId: mitra.ID,
			Gambar: os.Getenv("CLOUD_STORAGE_READ_LINK")+"galeriMitra/"+fileName,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		addedGaleri, err := s.galeriMitraRepository.Save(ctx, newGaleriMitra)

		if err != nil{
			return nil, err
		}

		newGaleriMitras = append(newGaleriMitras, addedGaleri.InsertedID.(primitive.ObjectID))
	}

	newMitra := bson.M{
		"galerimitra" : newGaleriMitras,
		"updatedat" : time.Now(),
	}

	updatedMitra, err := s.mitraRepository.UpdateProfil(ctx, mitra.ID, newMitra)

	if err != nil{
		return nil, err
	}

	return updatedMitra, nil
}