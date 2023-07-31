package services

import (
	"context"
	"errors"
	"fmt"
	"os"
	"pronics-api/constants"
	"pronics-api/formatters"
	"pronics-api/helper"
	"pronics-api/inputs"
	"pronics-api/models"
	"pronics-api/repositories"
	"sort"
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
	ShowKatalogMitra(ctx context.Context, userId primitive.ObjectID, searchFilter map[string] string) ([]formatters.KatalogResponse, error)
	ActivateMitra(ctx context.Context, mitraId primitive.ObjectID) (*mongo.UpdateResult, error)
	GetDetailMitra(ctx context.Context, mitraId primitive.ObjectID) (formatters.DetailMitraResponse, error)
	GetAllMitra(ctx context.Context) ([]formatters.MitraDashboardSummaryResponse, error)
	GetDashboardSummary(ctx context.Context, userId primitive.ObjectID) (formatters.DashboardSummaryMitra, error)
	GetAllLayananOwnedByMitra(ctx context.Context, mitraId primitive.ObjectID)([]formatters.LayananDetailMitraResponse, error)
	GetDetailMitraByAdmin(ctx context.Context, mitraId primitive.ObjectID) (formatters.DetailMitraByAdminResponse, error)
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
	komentarRepository repositories.KomentarRepository
	customerRepository repositories.CustomerRepository
	orderRepository repositories.OrderRepository
	orderDetailRepository repositories.OrderDetailRepository
	orderPaymentRepository repositories.OrderPaymentRepository
	ktpMitraRepository repositories.KTPMitraRepository
	savedRepository repositories.SavedRepository
}

func NewMitraService(userRepository repositories.UserRepository, mitraRepository repositories.MitraRepository, galeriMitraRepository repositories.GaleriRepository, wilayahRepository repositories.WilayahRepository, bidangRepository repositories.BidangRepository, kategoriRepository repositories.KategoriRepository, layananRepository repositories.LayananRepository, layananMitraRepository repositories.LayananMitraRepository, komentarRepository repositories.KomentarRepository, customerRepository repositories.CustomerRepository, orderRepository repositories.OrderRepository, orderDetailRepository repositories.OrderDetailRepository, orderPaymentRepository repositories.OrderPaymentRepository, ktpMitraRepository repositories.KTPMitraRepository, savedRepository repositories.SavedRepository) *mitraService{
	return &mitraService{userRepository, mitraRepository, galeriMitraRepository, wilayahRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, komentarRepository, customerRepository, orderRepository, orderDetailRepository, orderPaymentRepository, ktpMitraRepository, savedRepository}
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

	if err != nil{
		return detailBidang, err
	}

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

func (s *mitraService) ShowKatalogMitra(ctx context.Context, userId primitive.ObjectID, searchFilter map[string] string) ([]formatters.KatalogResponse, error){
	var katalogMitraResponses []formatters.KatalogResponse

	allMitra, err := s.mitraRepository.FindAllActiveMitra(ctx)

	if err != nil{
		return nil, err
	}

	textToSearch := strings.ToLower(searchFilter["search"])
	wilayahToSearch := strings.ToLower(searchFilter["daerah"])
	bidangToSearch := strings.ToLower(searchFilter["bidang"])
	sortBasedOn := strings.ToLower(searchFilter["urut"])
	alamatCustomer := strings.ToLower(searchFilter["alamatCustomer"])

	if alamatCustomer == ""{
		return nil, errors.New("alamat customer wajib diisi")
	}

	customer, err := s.customerRepository.GetCustomerByIdUser(ctx, userId)

	if err != nil{
		return nil, errors.New("customer not found")
	}

	for _, mitra := range allMitra{
		var katalogMitra formatters.KatalogResponse

		user, err := s.userRepository.GetUserById(ctx, mitra.UserId)

		if err != nil{
			return nil, err
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
				return nil, err
			}

			if(wilayahToSearch != strings.ToLower(wilayah.NamaWilayah)){
				continue
			}
		}


		katalogMitra.ID = mitra.ID
		katalogMitra.Gambar = mitra.GambarMitra

		comments, err := s.komentarRepository.GetAllByMitraId(ctx, mitra.ID)

		if err != nil {
			katalogMitra.Rating = 0
		}

		if len(comments) < 1{
			katalogMitra.Rating = 0
		}else{
			var sumRating float64
			for _, comment := range comments{
				sumRating += comment.Rating
			}
			katalogMitra.Rating = sumRating / float64(len(comments))
		}

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

		distance, err := helper.DistanceCalculation(alamatCustomer, mitra.Alamat)

		if err != nil{
			return nil, err
		}

		katalogMitra.Distance = distance

		_, err = s.savedRepository.GetByIdCustomerNMitra(ctx,customer.ID, mitra.ID)

		if err != nil{
			katalogMitra.IsSaved = false
		}else{
			katalogMitra.IsSaved = true
		}

		katalogMitraResponses = append(katalogMitraResponses, katalogMitra)
	}

	if(sortBasedOn != ""){
		if(sortBasedOn == constants.RatingTertinggi){
			sort.SliceStable(katalogMitraResponses, func(i, j int) bool {
				return katalogMitraResponses[i].Rating > katalogMitraResponses[j].Rating
			})
		}else if (sortBasedOn == constants.Terdekat){
			sort.SliceStable(katalogMitraResponses, func(i, j int) bool {
				return katalogMitraResponses[i].Distance < katalogMitraResponses[j].Distance
			})
		}else if(sortBasedOn == constants.Termurah){
			sort.SliceStable(katalogMitraResponses, func(i, j int) bool {
				return katalogMitraResponses[i].MinPrice < katalogMitraResponses[j].MinPrice
			})
		}
	}

	return katalogMitraResponses, nil
}

func (s *mitraService) ActivateMitra(ctx context.Context, mitraId primitive.ObjectID) (*mongo.UpdateResult, error){

	
	newMitra := bson.M{
		"status" : constants.MitraActive,
		"updatedat" : time.Now(),
	}

	updatedResult, err := s.mitraRepository.UpdateProfil(ctx, mitraId, newMitra)

	if err != nil{
		return nil, err
	}

	return updatedResult, nil
}

func (s *mitraService) GetDetailMitra(ctx context.Context, mitraId primitive.ObjectID) (formatters.DetailMitraResponse, error){

	var detailMitra formatters.DetailMitraResponse

	mitra, err := s.mitraRepository.GetMitraById(ctx, mitraId)

	if err != nil{
		return detailMitra, err
	}

	user, err := s.userRepository.GetUserById(ctx,mitra.UserId)

	if err != nil{
		return detailMitra, err
	}

	galeriImage, err := s.galeriMitraRepository.GetAllByIdMitra(ctx, mitra.ID)

	if err != nil{
		return detailMitra, errors.New("galeri image not found")
	}

	var galeriMitra []string

	for _, galeri := range galeriImage{
		galeriMitra = append(galeriMitra, galeri.Gambar)
	}

	var bidangArr []string
	var layananArr []formatters.LayananDetailMitraResponse

	for _, bidangId := range mitra.Bidang{
		bidangMitra, err := s.bidangRepository.GetById(ctx, bidangId)

		if err != nil{
			return detailMitra, errors.New("bidang mitra not found")
		}

		bidangArr = append(bidangArr, bidangMitra.NamaBidang)

		layanan, err := s.layananRepository.FindAllByBidangId(ctx, bidangMitra.ID)

		if err != nil{
			return detailMitra, errors.New("layanan not found")
		}

		layananMitra, err := s.layananMitraRepository.FindAllByBidangAndMitra(ctx, bidangMitra.ID, mitra.ID)

		if err != nil{
			return detailMitra, errors.New("layanan mitra not found")
		}

		for _, item := range layanan{
			var layananForResponse formatters.LayananDetailMitraResponse

			layananForResponse.ID = item.ID
			layananForResponse.NamaLayanan = item.NamaLayanan
			layananForResponse.Harga = item.Harga
			layananForResponse.BidangId = bidangMitra.ID

			layananArr = append(layananArr, layananForResponse)
		}

		for _, item := range layananMitra{
			var layananForResponse formatters.LayananDetailMitraResponse

			layananForResponse.ID = item.ID
			layananForResponse.NamaLayanan = item.NamaLayanan
			layananForResponse.Harga = item.Harga
			layananForResponse.BidangId = bidangMitra.ID

			layananArr = append(layananArr, layananForResponse)
		}
	}

	comments, err := s.komentarRepository.GetAllByMitraId(ctx, mitra.ID)

	if err != nil{
		detailMitra.Ulasan = nil
	}

	var averageRating = 0.0
	var commentsResponse formatters.KomentarDetailMitraResponse
	
	for _, comment := range comments{
		var commentResponse formatters.KomentarResponse

		customer, err := s.customerRepository.GetCustomerById(ctx, comment.CustomerId)

		if err != nil{
			return detailMitra, err
		}

		order, err := s.orderRepository.GetById(ctx, comment.OrderId)

		if err != nil{
			return detailMitra, err
		}

		orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

		if err != nil{
			return detailMitra, err
		}

		layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)
		var namaLayanan string

		if err != nil{
			layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

			if err != nil{
				return detailMitra, err
			}

			namaLayanan = layananMitra.NamaLayanan
		}else{
			namaLayanan = layanan.NamaLayanan
		}

		user, err := s.userRepository.GetUserById(ctx, customer.UserId)

		if err != nil{
			return detailMitra, err
		}

		commentResponse.ID = comment.ID
		commentResponse.FotoCustomer = customer.GambarCustomer
		commentResponse.Gambar = comment.GambarKomentar
		commentResponse.Komentar = comment.Komentar
		commentResponse.Layanan = namaLayanan
		commentResponse.NamaCustomer = user.NamaLengkap
		commentResponse.RatingGiven = comment.Rating
		commentResponse.Tanggal = comment.UpdatedAt
		commentResponse.TotalSuka = len(comment.Penyuka)

		averageRating = averageRating + comment.Rating
		
		commentsResponse.AllKomentar = append(commentsResponse.AllKomentar, commentResponse)
	}

	banyakComment := len(comments)

	commentsResponse.UlasanCount = banyakComment

	if banyakComment == 0{
		commentsResponse.OverallRating = 0.0
	}else{
		commentsResponse.OverallRating = averageRating / float64(banyakComment)
	}


	detailMitra.ID = mitra.ID
	detailMitra.NamaPemilik = user.NamaLengkap
	detailMitra.NamaToko = mitra.NamaToko
	detailMitra.GaleriImage = galeriMitra
	detailMitra.FotoProfil = mitra.GambarMitra
	detailMitra.Deskripsi = user.Deskripsi
	detailMitra.Bidang = bidangArr
	detailMitra.Layanan = layananArr
	detailMitra.Ulasan = commentsResponse
	

	return detailMitra, nil

}

func (s *mitraService) GetAllMitra(ctx context.Context) ([]formatters.MitraDashboardSummaryResponse, error){
	var mitraResponses []formatters.MitraDashboardSummaryResponse

	mitras, err := s.mitraRepository.GetAllMitra(ctx)

	if err != nil{
		return mitraResponses, err
	}

	for _, mitra := range mitras{
		var mitraResponse formatters.MitraDashboardSummaryResponse

		user, err := s.userRepository.GetUserById(ctx, mitra.UserId)

		if err != nil{
			return mitraResponses, err
		}

		orders, err := s.orderRepository.GetAllOrderMitraSelesai(ctx, mitra.ID, constants.OrderCompleted)

		if err != nil{
			mitraResponse.JumlahTransaksiSelesai = 0
		}else{
			mitraResponse.JumlahTransaksiSelesai = len(orders)
		}

		mitraResponse.ID = mitra.ID
		mitraResponse.Email = user.Email
		mitraResponse.NamaPemilik = user.NamaLengkap
		mitraResponse.NamaToko = mitra.NamaToko
		mitraResponse.NoHandphone = user.NoTelepon
		mitraResponse.StatusMitra = mitra.Status

		mitraResponses = append(mitraResponses, mitraResponse)
	}

	return mitraResponses, nil
}

func (s *mitraService) GetDashboardSummary(ctx context.Context, userId primitive.ObjectID) (formatters.DashboardSummaryMitra, error){
	var dashboardSummaryMira formatters.DashboardSummaryMitra

	mitra, err := s.mitraRepository.GetMitraByIdUser(ctx, userId)

	if err != nil{
		return dashboardSummaryMira, errors.New("mitra not found")
	}

	orders, err := s.orderRepository.GetAllOrderMitraSelesai(ctx, mitra.ID, constants.OrderCompleted)

	var totalOrderSelesai int
	var totalPendapatanBersih float64
	
	if err != nil{
		totalOrderSelesai = 0
		totalPendapatanBersih = 0.0
	}else{
		for _, order := range orders{
			orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

			if err != nil{
				return dashboardSummaryMira, err
			}

			orderPayment, err := s.orderPaymentRepository.GetByOrderDetailId(ctx, orderDetail.ID)

			if err != nil{
				return dashboardSummaryMira, err
			}

			totalOrderSelesai += 1

			pendapatanBersih := orderPayment.TotalBiaya - orderPayment.BiayaAplikasi

			totalPendapatanBersih += pendapatanBersih
		}
	}

	dashboardSummaryMira.TotalOrderSelesai = totalOrderSelesai
	dashboardSummaryMira.TotalPendapatanBersih = totalPendapatanBersih

	return dashboardSummaryMira, nil
}

func (s *mitraService) GetAllLayananOwnedByMitra(ctx context.Context, mitraId primitive.ObjectID)([]formatters.LayananDetailMitraResponse, error){
	var layanansOwnedByMitra []formatters.LayananDetailMitraResponse

	mitra, err := s.mitraRepository.GetMitraById(ctx, mitraId)

	if err != nil{
		return layanansOwnedByMitra, errors.New("mitra not found")
	}

	for _, bidangId := range mitra.Bidang{
		bidang, err := s.bidangRepository.GetById(ctx, bidangId)

		if err != nil{
			return layanansOwnedByMitra, err
		}

		for _, layananId := range bidang.LayananId{
			var layananOwnedByMitra formatters.LayananDetailMitraResponse
			layanan, err := s.layananRepository.GetById(ctx, layananId)

			if err != nil{
				layananMitra, err := s.layananMitraRepository.GetById(ctx, layananId)

				if err != nil{
					return layanansOwnedByMitra, err
				}

				layananOwnedByMitra.ID = layananMitra.ID
				layananOwnedByMitra.NamaLayanan = layananMitra.NamaLayanan
				layananOwnedByMitra.Harga = layananMitra.Harga
				layananOwnedByMitra.BidangId = bidangId

				
			}else{
				layananOwnedByMitra.ID = layanan.ID
				layananOwnedByMitra.NamaLayanan = layanan.NamaLayanan
				layananOwnedByMitra.Harga = layanan.Harga
				layananOwnedByMitra.BidangId = bidangId
			}

			layanansOwnedByMitra = append(layanansOwnedByMitra, layananOwnedByMitra)
		}
	}

	return layanansOwnedByMitra, nil
}

func (s *mitraService) GetDetailMitraByAdmin(ctx context.Context, mitraId primitive.ObjectID) (formatters.DetailMitraByAdminResponse, error){
	var detailMitraByAdmin formatters.DetailMitraByAdminResponse
	var detailMitra formatters.DetailMitraResponse

	mitra, err := s.mitraRepository.GetMitraById(ctx, mitraId)

	if err != nil{
		return detailMitraByAdmin, err
	}

	user, err := s.userRepository.GetUserById(ctx,mitra.UserId)

	if err != nil{
		return detailMitraByAdmin, err
	}

	galeriImage, err := s.galeriMitraRepository.GetAllByIdMitra(ctx, mitra.ID)

	if err != nil{
		return detailMitraByAdmin, err
	}

	var galeriMitra []string

	for _, galeri := range galeriImage{
		galeriMitra = append(galeriMitra, galeri.Gambar)
	}

	var bidangArr []string
	var layananArr []formatters.LayananDetailMitraResponse

	for _, bidangId := range mitra.Bidang{
		bidangMitra, err := s.bidangRepository.GetById(ctx, bidangId)

		if err != nil{
			return detailMitraByAdmin, err
		}

		bidangArr = append(bidangArr, bidangMitra.NamaBidang)

		layanan, err := s.layananRepository.FindAllByBidangId(ctx, bidangMitra.ID)

		if err != nil{
			return detailMitraByAdmin, err
		}

		layananMitra, err := s.layananMitraRepository.FindAllByBidangAndMitra(ctx, bidangMitra.ID, mitra.ID)

		if err != nil{
			return detailMitraByAdmin, err
		}

		for _, item := range layanan{
			var layananForResponse formatters.LayananDetailMitraResponse

			layananForResponse.ID = item.ID
			layananForResponse.NamaLayanan = item.NamaLayanan
			layananForResponse.Harga = item.Harga
			layananForResponse.BidangId = bidangMitra.ID

			layananArr = append(layananArr, layananForResponse)
		}

		for _, item := range layananMitra{
			var layananForResponse formatters.LayananDetailMitraResponse

			layananForResponse.ID = item.ID
			layananForResponse.NamaLayanan = item.NamaLayanan
			layananForResponse.Harga = item.Harga
			layananForResponse.BidangId = bidangMitra.ID

			layananArr = append(layananArr, layananForResponse)
		}
	}

	comments, err := s.komentarRepository.GetAllByMitraId(ctx, mitra.ID)

	if err != nil{
		detailMitra.Ulasan = nil
	}

	var averageRating float64
	var commentsResponse formatters.KomentarDetailMitraResponse
	
	for _, comment := range comments{
		var commentResponse formatters.KomentarResponse

		customer, err := s.customerRepository.GetCustomerById(ctx, comment.CustomerId)

		if err != nil{
			return detailMitraByAdmin, err
		}

		order, err := s.orderRepository.GetById(ctx, comment.OrderId)

		if err != nil{
			return detailMitraByAdmin, err
		}

		orderDetail, err := s.orderDetailRepository.GetByOrderId(ctx, order.ID)

		if err != nil{
			return detailMitraByAdmin, err
		}

		layanan, err := s.layananRepository.GetById(ctx, orderDetail.LayananId)
		var namaLayanan string

		if err != nil{
			layananMitra, err := s.layananMitraRepository.GetById(ctx, orderDetail.LayananId)

			if err != nil{
				return detailMitraByAdmin, err
			}

			namaLayanan = layananMitra.NamaLayanan
		}else{
			namaLayanan = layanan.NamaLayanan
		}

		user, err := s.userRepository.GetUserById(ctx, customer.UserId)

		if err != nil{
			return detailMitraByAdmin, err
		}

		commentResponse.ID = comment.ID
		commentResponse.FotoCustomer = customer.GambarCustomer
		commentResponse.Gambar = comment.GambarKomentar
		commentResponse.Komentar = comment.Komentar
		commentResponse.Layanan = namaLayanan
		commentResponse.NamaCustomer = user.NamaLengkap
		commentResponse.RatingGiven = comment.Rating
		commentResponse.Tanggal = comment.UpdatedAt
		commentResponse.TotalSuka = len(comment.Penyuka)

		averageRating = averageRating + comment.Rating
		
		commentsResponse.AllKomentar = append(commentsResponse.AllKomentar, commentResponse)
	}

	commentsResponse.UlasanCount = len(comments)
	commentsResponse.OverallRating = averageRating / float64(len(comments))



	detailMitra.ID = mitra.ID
	detailMitra.NamaPemilik = user.NamaLengkap
	detailMitra.NamaToko = mitra.NamaToko
	detailMitra.GaleriImage = galeriMitra
	detailMitra.FotoProfil = mitra.GambarMitra
	detailMitra.Deskripsi = user.Deskripsi
	detailMitra.Bidang = bidangArr
	detailMitra.Layanan = layananArr
	detailMitra.Ulasan = commentsResponse

	detailMitraByAdmin.DetailMitra = detailMitra

	ktp, err := s.ktpMitraRepository.GetByMitraId(ctx, mitraId)

	if err != nil{
		return detailMitraByAdmin, errors.New("ktp is not found")
	}

	detailMitraByAdmin.KTP = ktp.GambarKtp

	return detailMitraByAdmin, nil
}