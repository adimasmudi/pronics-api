package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func MitraRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection, galeriMitraCollection *mongo.Collection, wilayahCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection, komentarCollection *mongo.Collection,customerCollection *mongo.Collection, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection,orderPaymentCollection *mongo.Collection, adminCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	galeriMitraRepository := repositories.NewGaleriRepository(galeriMitraCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	wilayahRepository := repositories.NewWilayahRepository(wilayahCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)
	komentarRepository := repositories.NewKomentarRepository(komentarCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	orderDetailRepository := repositories.NewOrderDetailRepository(orderDetailCollection)
	orderPaymentRepository := repositories.NewOrderPaymentRepository(orderPaymentCollection)
	adminRepository := repositories.NewAdminRepository(adminCollection)
	
	// services
	mitraService := services.NewMitraService(userRepository, mitraRepository, galeriMitraRepository, wilayahRepository,bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, komentarRepository, customerRepository, orderRepository, orderDetailRepository, orderPaymentRepository)

	// controllers
	mitraHandler := controllers.NewMitraHandler(mitraService)

	// auth
	adminAuth := middlewares.AdminAuth(adminRepository)
	// auth
	mitra := middlewares.MitraAuth(mitraRepository)

	mitraRoute := api.Group("/mitra")

	mitraRoute.Get("/profile", mitra.AuthMitra, mitraHandler.GetProfile)
	mitraRoute.Put("/profile/update", mitra.AuthMitra, mitraHandler.UpdateProfile)
	mitraRoute.Put("/galeri/upload", mitra.AuthMitra, mitraHandler.UploadMultipleImagesToGaleri)
	mitraRoute.Put("/updateBidang", mitra.AuthMitra, mitraHandler.UpdateBidang)
	mitraRoute.Get("/getBidangs",mitra.AuthMitra, mitraHandler.GetBidangMitra)
	mitraRoute.Get("/getBidang/detail/:bidangId", mitra.AuthMitra, mitraHandler.DetailBidangMitra)
	mitraRoute.Get("/showKatalog", mitraHandler.ShowKatalogMitra)
	mitraRoute.Put("/activate/:mitraId", adminAuth.AuthAdmin, mitraHandler.ActivateMitra)
	mitraRoute.Get("/detail/:mitraId", mitraHandler.DetailMitra)
	mitraRoute.Get("/all", adminAuth.AuthAdmin, mitraHandler.GetAllMitra)
	mitraRoute.Get("/dashboardSummary", mitra.AuthMitra, mitraHandler.GetDashboardSummary)
}