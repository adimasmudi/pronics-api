package services

import (
	"context"
	"errors"
	"os"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/exp/slices"
)

type KomentarService interface {
	AddKomentar(ctx context.Context, orderId primitive.ObjectID, input inputs.KomentarInput, fileNames []string) (*mongo.InsertOneResult, error)
	SeeKomentar(ctx context.Context, orderId primitive.ObjectID) (formatters.KomentarResponse, error)
	UpdateKomentar(ctx context.Context, komentarId primitive.ObjectID, input inputs.KomentarInput, fileNames []string) (*mongo.UpdateResult, error)
	ResponseKomentar(ctx context.Context, userId primitive.ObjectID, komentarId primitive.ObjectID, tipe string) (*mongo.UpdateResult, error)
	// DeleteKomentar(ctx context.Context, komentarId primitive.ObjectID) (*mongo.DeleteResult, error)
}

type komentarService struct{
	userRepository     repositories.UserRepository
	mitraRepository repositories.MitraRepository
	customerRepository repositories.CustomerRepository
	orderRepository repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	komentarRepository repositories.KomentarRepository
	layananRepository repositories.LayananRepository
	layananMitraRepository repositories.LayananMitraRepository
}

func NewKomentarService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, komentarRepository repositories.KomentarRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository) *komentarService{
	return &komentarService{userRepository, mitraRepository, customerRepository, orderRepository, orderDetailRepository, komentarRepository, layananRepository, layananMitraRepository}
}

// add komentar
func (s *komentarService) AddKomentar(ctx context.Context,orderId primitive.ObjectID, input inputs.KomentarInput, fileNames []string) (*mongo.InsertOneResult, error){
	
	order, err := s.orderRepository.GetById(ctx, orderId)

	if err != nil{
		return nil, errors.New("order with that id not found")
	}

	customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

	if err != nil{
		return nil, err
	}

	mitra, err := s.mitraRepository.GetMitraById(ctx, order.MitraId)

	if err != nil{
		return nil, err
	}

	var gambarKomentarArr []string
	for _, fileName := range fileNames{
		fileName = os.Getenv("CLOUD_STORAGE_READ_LINK")+"komentar/"+fileName
		gambarKomentarArr = append(gambarKomentarArr, fileName)
	}

	newKomentar := models.Komentar{
		ID : primitive.NewObjectID(),
		CustomerId: customer.ID,
		MitraId: mitra.ID,
		OrderId: order.ID,
		Rating: input.Rating,
		Komentar: input.Komentar,
		GambarKomentar: gambarKomentarArr,
		Penyuka: []primitive.ObjectID{},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	addKomentar, err := s.komentarRepository.Save(ctx, newKomentar)

	if err != nil{
		return nil, err
	}

	return addKomentar, nil
}

// see komentar
func (s *komentarService) SeeKomentar(ctx context.Context, orderId primitive.ObjectID) (formatters.KomentarResponse, error){
	var komentarResponse formatters.KomentarResponse

	order, err := s.orderRepository.GetById(ctx, orderId)

	if err != nil{
		return komentarResponse, errors.New("order with that id not found")
	}

	customer, err := s.customerRepository.GetCustomerById(ctx, order.CustomerId)

	if err != nil{
		return komentarResponse, err
	}

	user, err := s.userRepository.GetUserById(ctx, customer.UserId)

	if err != nil{
		return komentarResponse, err
	}

	orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

	if err != nil{
		return komentarResponse, err
	}

	komentar, err := s.komentarRepository.GetByOrderId(ctx, order.ID)

	if err != nil{
		return komentarResponse, err
	}

	layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)
	var namaLayanan string

	if err != nil{
		layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

		if err != nil{
			return komentarResponse, err
		}

		namaLayanan = layananMitra.NamaLayanan
	}else{
		namaLayanan = layanan.NamaLayanan
	}

	komentarResponse.ID = komentar.ID
	komentarResponse.FotoCustomer = customer.GambarCustomer
	komentarResponse.Gambar = komentar.GambarKomentar
	komentarResponse.Komentar = komentar.Komentar
	komentarResponse.Layanan = namaLayanan
	komentarResponse.NamaCustomer = user.NamaLengkap
	komentarResponse.RatingGiven = komentar.Rating
	komentarResponse.Tanggal = komentar.UpdatedAt
	komentarResponse.TotalSuka = len(komentar.Penyuka)

	return komentarResponse, nil

}

// update komentar
func (s *komentarService) UpdateKomentar(ctx context.Context,komentarId primitive.ObjectID, input inputs.KomentarInput, fileNames []string) (*mongo.UpdateResult, error){
	
	komentar, err := s.komentarRepository.GetById(ctx, komentarId)

	if err != nil{
		return nil, err
	}

	gambarKomentarArr := komentar.GambarKomentar

	for _, fileName := range fileNames{
		fileName = os.Getenv("CLOUD_STORAGE_READ_LINK")+"komentar/"+fileName
		gambarKomentarArr = append(gambarKomentarArr, fileName)
	}

	newKomentar := bson.M{
		"rating":  input.Rating,
		"komentar" : input.Komentar,
		"gambar_komentar" : gambarKomentarArr,
		"updatedat" : time.Now(),
	}

	updatedKomentar, err := s.komentarRepository.Update(ctx, komentarId,newKomentar)

	if err != nil{
		return nil, err
	}

	return updatedKomentar, nil
}

// response komentar
func (s *komentarService) ResponseKomentar(ctx context.Context, userId primitive.ObjectID, komentarId primitive.ObjectID, tipe string) (*mongo.UpdateResult, error){
	if(tipe == "" || (tipe != constants.LikeComment && tipe != constants.UnLikeComment)){
		return nil, errors.New("tipe wajib diisi dengan nilai antara 'like' atau 'unlike'")
	}
	
	user, err := s.userRepository.GetUserById(ctx, userId)

	if err != nil{
		return nil, err
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, user.ID)

	if err != nil{
		return nil, errors.New("customer not found")
	}

	komentar, err := s.komentarRepository.GetById(ctx, komentarId)

	if err != nil{
		return nil, errors.New("error select komentar id")
	}

	penyuka := komentar.Penyuka

	if(tipe == constants.LikeComment){
		for _, suka := range penyuka{
			if suka == customer.ID{
				return nil, errors.New("kamu sudah like komentar ini")
			}
		}
	
		penyuka = append(penyuka, customer.ID)
	}else if (tipe == constants.UnLikeComment){
		isAlreadyLike := false
		for idx, suka := range penyuka{
			if suka == customer.ID{
				penyuka = slices.Delete(penyuka, idx, idx+1)
				isAlreadyLike = true
			}
		}

		if !isAlreadyLike{
			return nil, errors.New("kamu belum like komentar ini")
		}
	}

	newKomentar := bson.M{
		"penyuka" : penyuka,
		"updatedat" : time.Now(),
	}

	updatedKomentar, err := s.komentarRepository.Update(ctx, komentarId,newKomentar)

	if err != nil{
		return nil, err
	}

	return updatedKomentar, nil
}

// delete komentar