package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func BidangRoute(api fiber.Router, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, adminCollection *mongo.Collection){
	// repositories
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	adminRepository := repositories.NewAdminRepository(adminCollection)

	// services
	bidangService := services.NewbidangService(bidangRepository, kategoriRepository, layananRepository)

	// controllers
	bidangHandler := controllers.NewbidangHandler(bidangService)

	// auth
	adminAuth := middlewares.AdminAuth(adminRepository)

	bidang := api.Group("/bidang")

	bidang.Post("/save", adminAuth.AuthAdmin, bidangHandler.Save)
	bidang.Get("/all", bidangHandler.FindAll)
	bidang.Get("/detail/:bidangId", bidangHandler.FindById)
	bidang.Put("/update/:bidangId", adminAuth.AuthAdmin, bidangHandler.UpdateBidang)
	bidang.Delete("/delete/:bidangId", adminAuth.AuthAdmin, bidangHandler.DeleteBidang)
}