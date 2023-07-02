package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/models"
	"pronics-api/repositories"
	"sort"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type SavedService interface {
	Save(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (*mongo.InsertOneResult, error)
	ShowAll(ctx context.Context, userId primitive.ObjectID, searchFilter map[string] string)([]formatters.SavedResponse, error)
	DeleteSaved(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type savedService struct{
	userRepository     repositories.UserRepository
	customerRepository repositories.CustomerRepository
	mitraRepository repositories.MitraRepository
	wilayahRepository repositories.WilayahRepository
	bidangRepository repositories.BidangRepository
	kategoriRepository repositories.KategoriRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
	savedRepository repositories.SavedRepository
}

func NewSavedService(userRepository repositories.UserRepository,customerRepository repositories.CustomerRepository, mitraRepository repositories.MitraRepository,wilayahRepository repositories.WilayahRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository, savedRepository repositories.SavedRepository) *savedService{
	return &savedService{userRepository, customerRepository,mitraRepository,wilayahRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, savedRepository}
}

func (s *savedService) Save(ctx context.Context, userId primitive.ObjectID, mitraId primitive.ObjectID) (*mongo.InsertOneResult, error){
	
	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return nil, err
	}

	checkSavedExist, err := s.savedRepository.GetByIdMitra(ctx,mitraId)

	if err == nil{
		fmt.Println(checkSavedExist)
		return nil, errors.New("anda sudah menyimpan mitra tersebut sebelumnya")
	}

	newSaved := models.Saved{
		ID : primitive.NewObjectID(),
		Customer_id: customer.ID,
		Mitra_id: mitraId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	savedAdded, err := s.savedRepository.Save(ctx, newSaved)

	if err != nil{
		return nil, err
	}

	return savedAdded, nil
}

func (s *savedService) ShowAll(ctx context.Context, userId primitive.ObjectID, searchFilter map[string] string)([]formatters.SavedResponse, error){
	var savedMitras []formatters.SavedResponse

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return savedMitras, errors.New("customer not found")
	}

	textToSearch := strings.ToLower(searchFilter["search"])
	wilayahToSearch := strings.ToLower(searchFilter["daerah"])
	bidangToSearch := strings.ToLower(searchFilter["bidang"])
	sortBasedOn := strings.ToLower(searchFilter["urut"])
	alamatCustomer := strings.ToLower(searchFilter["alamatCustomer"])

	if alamatCustomer == ""{
		return nil, errors.New("alamat customer wajib diisi")
	}
	
	allSavedMitra, err := s.savedRepository.GetAll(ctx,customer.ID)

	if err != nil{
		return savedMitras, errors.New("saved mitra not found")
	}

	for _, savedMitra := range allSavedMitra{
		var savedResponse formatters.SavedResponse
		var katalogMitra formatters.KatalogResponse
		mitra, err := s.mitraRepository.GetMitraById(ctx, savedMitra.Mitra_id)

		if err != nil{
			return savedMitras, errors.New("mitra not found")
		}

		user, err := s.userRepository.GetUserById(ctx, mitra.UserId)

		if err != nil{
			return nil, errors.New("mitra user information not found")
		}

		if(mitra.NamaToko != ""){
			katalogMitra.Name = mitra.NamaToko
		}else{
			katalogMitra.Name = user.NamaLengkap
		}

		if(textToSearch != ""){
			if(!(strings.Contains(strings.ToLower(katalogMitra.Name), textToSearch) || strings.Contains(strings.ToLower(user.Deskripsi), textToSearch))){
				continue
			}
		}

		if(wilayahToSearch != ""){
			wilayah, err := s.wilayahRepository.FindById(ctx, mitra.Wilayah)

			if err != nil{
				return nil, errors.New("wilayah mitra not found")
			}

			if(wilayahToSearch != strings.ToLower(wilayah.NamaWilayah)){
				continue
			}
		}

		katalogMitra.ID = mitra.ID
		katalogMitra.Gambar = mitra.GambarMitra
		katalogMitra.Rating = 5 // sementara

		min := 0.0
		max := 0.0

		for _, bidangId := range mitra.Bidang{
			var bidang formatters.BidangResponse
			bidangMitra, err := s.bidangRepository.GetById(ctx, bidangId)

			if err != nil{
				return nil, err
			}

			bidang.ID = bidangMitra.ID

			kategori, err := s.kategoriRepository.GetById(ctx, bidangMitra.KategoriId)

			if err != nil{
				return nil, err
			}
			bidang.Kategori = kategori.NamaKategori
			bidang.NamaBidang = bidangMitra.NamaBidang

			katalogMitra.Bidang = append(katalogMitra.Bidang, bidang)

			layanan, err := s.layananRepository.FindAllByBidangId(ctx, bidangMitra.ID)

			if err != nil{
				return nil, err
			}

			layananMitra, err := s.layananMitraRepository.FindAllByBidangAndMitra(ctx, bidangMitra.ID, mitra.ID)

			if err != nil{
				return nil, err
			}

			for _, item := range layanan{
				if min == 0{
					min = item.Harga
				}else if min > item.Harga{
					min = item.Harga
				}

				if max == 0{
					max = item.Harga
				}else if max < item.Harga{
					max = item.Harga
				}
			}

			for _, item := range layananMitra{
				if min == 0{
					min = item.Harga
				}else if min > item.Harga{
					min = item.Harga
				}

				if max == 0{
					max = item.Harga
				}else if max < item.Harga{
					max = item.Harga
				}
			}

		}

		bidangContains := false
		if(bidangToSearch != ""){
			for _, bidang := range katalogMitra.Bidang{
				if strings.ToLower(bidang.NamaBidang) == bidangToSearch{
					bidangContains = true
				}
			}

			if !bidangContains{
				continue
			}
		}

		katalogMitra.MinPrice = min
		katalogMitra.MaxPrice = max

		url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&units=imperial&key=%s",alamatCustomer,mitra.Alamat,os.Getenv("MAPS_API_KEY"))
		url = strings.Replace(url, " ", "%20", -1)
		method := "GET"

		fmt.Println(url)

		client := &http.Client {}
		req, err := http.NewRequest(method, url, nil)

		if err != nil {
			return nil, err
		}
		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		
		defer res.Body.Close()

		body, err := io.ReadAll(res.Body)

		if err != nil {
			return nil, err
		}

		var distanceMatrix helper.DistanceMatrixResult

		
		err = json.Unmarshal(body, &distanceMatrix)
		if err != nil {
			return nil, err
		}

		fmt.Println(distanceMatrix)

		katalogMitra.Distance = distanceMatrix.Rows[0].Elements[0].Distance.Value

		savedResponse.ID = savedMitra.ID
		savedResponse.Mitra = katalogMitra

		savedMitras = append(savedMitras, savedResponse)
	}

	if(sortBasedOn != ""){
		if(sortBasedOn == constants.RatingTertinggi){
			sort.SliceStable(savedMitras, func(i, j int) bool {
				return savedMitras[i].Mitra.Rating > savedMitras[j].Mitra.Rating
			})
		}else if (sortBasedOn == constants.Terdekat){
			sort.SliceStable(savedMitras, func(i, j int) bool {
				return savedMitras[i].Mitra.Distance < savedMitras[j].Mitra.Distance
			})
		}else if(sortBasedOn == constants.Termurah){
			sort.SliceStable(savedMitras, func(i, j int) bool {
				return savedMitras[i].Mitra.MinPrice < savedMitras[j].Mitra.MinPrice
			})
		}
	}

	return savedMitras, nil


}

func (s *savedService) DeleteSaved(ctx context.Context, savedId primitive.ObjectID) (*mongo.DeleteResult, error){
	deletedSaved, err := s.savedRepository.Delete(ctx, savedId)

	if err != nil{
		return deletedSaved, err
	}

	return deletedSaved, nil
}