package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func LayananRoute(api fiber.Router, layananCollection *mongo.Collection, bidangCollection *mongo.Collection, adminCollection *mongo.Collection){
	// repositories
	adminRepository := repositories.NewAdminRepository(adminCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)

	// services
	layananService := services.NewLayananService(layananRepository, bidangRepository)

	// controllers
	layananHandler := controllers.NewLayananHandler(layananService)

	// auth
	adminAuth := middlewares.AdminAuth(adminRepository)

	layanan := api.Group("/layanan")

	layanan.Post("/save", adminAuth.AuthAdmin, layananHandler.Save)
	layanan.Get("/all/byBidang/:bidangId",layananHandler.FindAllByBidang)
	layanan.Get("/detail/:layananId", layananHandler.FindById)
	layanan.Delete("/delete/:layananId", adminAuth.AuthAdmin, layananHandler.DeleteLayanan)
}