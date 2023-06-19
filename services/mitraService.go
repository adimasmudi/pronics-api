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
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type MitraService interface {
	GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error)
	UpdateProfileMitra(ctx context.Context, ID primitive.ObjectID, input inputs.UpdateProfilMitraInput, fileName string) (*mongo.UpdateResult, error)
	UploadGaleriImage(ctx context.Context, ID primitive.ObjectID, fileNames []string) (*mongo.UpdateResult, error)
	GetAllBidangMitra(ctx context.Context, userId primitive.ObjectID) ([]formatters.BidangResponse, error)
	UpdateBidang(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateBidangMitraInput) (*mongo.UpdateResult, error)
	DetailBidang(ctx context.Context,userId primitive.ObjectID, bidangId primitive.ObjectID) (formatters.DetailBidangResponse, error)
}

type mitraService struct {
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
	galeriMitraRepository repositories.GaleriRepository
	wilayahRepository repositories.WilayahRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewMitraService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, galeriMitraRepository repositories.GaleriRepository, wilayahRepository repositories.WilayahRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *mitraService{
	return &mitraService{userRepository, mitraRepository, galeriMitraRepository, wilayahRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository}
}

func (s *mitraService) GetMitraProfile(ctx context.Context, ID primitive.ObjectID) (formatters.MitraResponse, error){
	var data formatters.MitraResponse

	user, err := s.userRepository.GetUserById(ctx, ID)

	if err != nil{
		return data, err
	}

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, ID)

	if err != nil{
		return data, err
	}

	wilayahMitra, err := s.wilayahRepository.FindById(ctx, mitra.Wilayah)

	if err != nil{
		return data, err
	}

	var bidangs []models.Bidang

	for _, bidangIdMitra := range mitra.Bidang{
		bidangMitra, err := s.bidangRepository.GetById(ctx, bidangIdMitra)

		if err != nil{
			return data, err
		}

		bidangs = append(bidangs, bidangMitra)
	}

	data = helper.MapperMitra(user, mitra, wilayahMitra, bidangs)

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
		return nil, errors.New("tidak ada gambar di upload")
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

func (s *mitraService) GetAllBidangMitra(ctx context.Context, userId primitive.ObjectID) ([]formatters.BidangResponse, error){
	var allBidangMitras []formatters.BidangResponse

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}

	for _, bidangId := range mitra.Bidang{
		var bidangResponse formatters.BidangResponse

		bidang, err := s.bidangRepository.GetById(ctx, bidangId)

		if err != nil{
			return nil, err
		}

		kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

		if err != nil{
			return nil, err
		}

		bidangResponse.ID = bidang.ID
		bidangResponse.Kategori = kategori.NamaKategori
		bidangResponse.NamaBidang = bidang.NamaBidang

		allBidangMitras = append(allBidangMitras, bidangResponse)
	}

	return allBidangMitras, nil
}

func (s *mitraService) UpdateBidang(ctx context.Context, userId primitive.ObjectID, input inputs.UpdateBidangMitraInput) (*mongo.UpdateResult, error){
	var newMitra primitive.M

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}

	bidangStrArr := input.Bidang
	bidangStr := strings.TrimSpace(bidangStrArr)
	bidangStr = strings.Replace(bidangStr, "[", "", -1)
	bidangStr = strings.Replace(bidangStr, "]", "", -1)
	bidangArr := strings.Split(bidangStr, ",")

	var bidangMitra []primitive.ObjectID

	for _, bidang := range bidangArr{
		bidang = strings.Trim(bidang," ")
		eachBidang, _ := primitive.ObjectIDFromHex(bidang)
		bidangMitra = append(bidangMitra, eachBidang)
	}

	newMitra = bson.M{
		"bidang" : bidangMitra,
		"updatedat" : time.Now(),
	}

	updatedMitra, err := s.mitraRepository.UpdateProfil(ctx, mitra.ID, newMitra)

	if err != nil{
		return nil, err
	}

	return updatedMitra, nil
}

func (s *mitraService) DetailBidang(ctx context.Context,userId primitive.ObjectID,  bidangId primitive.ObjectID) (formatters.DetailBidangResponse, error){
	var detailBidang formatters.DetailBidangResponse

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	bidang, err := s.bidangRepository.GetById(ctx, bidangId)

	if err != nil{
		return detailBidang, err
	}

	kategori, err := s.kategoriRepository.GetById(ctx, bidang.KategoriId)

	if err != nil{
		return detailBidang, err
	}

	layananInBidang, err := s.layananRepository.FindAllByBidangId(ctx, bidang.ID)

	if err != nil{
		return detailBidang, err
	}

	layananMitraInBidang, err := s.layananMitraRepository.FindAllByBidangAndMitra(ctx, bidang.ID, mitra.ID)

	if err != nil{
		return detailBidang, err
	}

	var allLayanan []formatters.LayananResponse

	for _, layananItem := range layananInBidang{
		var layanan formatters.LayananResponse

		layanan.ID = layananItem.ID
		layanan.NamaLayanan = layananItem.NamaLayanan

		allLayanan = append(allLayanan, layanan)
	}

	for _, layananItem := range layananMitraInBidang{
		var layanan formatters.LayananResponse

		layanan.ID = layananItem.ID
		layanan.NamaLayanan = layananItem.NamaLayanan

		allLayanan = append(allLayanan, layanan)
	}

	detailBidang.ID = bidang.ID
	detailBidang.NamaBidang = bidang.NamaBidang
	detailBidang.Kategori = kategori.NamaKategori
	detailBidang.Layanan = allLayanan

	return detailBidang, nil
}