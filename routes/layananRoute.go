package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func LayananRoute(api fiber.Router, layananCollection *mongo.Collection, bidangCollection *mongo.Collection){
	// repositories
	layananRepository := repositories.NewLayananRepository(layananCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)

	// services
	layananService := services.NewLayananService(layananRepository, bidangRepository)

	// controllers
	layananHandler := controllers.NewLayananHandler(layananService)

	layanan := api.Group("/layanan")

	layanan.Post("/save", middlewares.Auth, layananHandler.Save)
	layanan.Get("/all/byBidang/:bidangId",layananHandler.FindAllByBidang)
	layanan.Get("/detail/:layananId", layananHandler.FindById)
	layanan.Delete("/delete/:layananId", middlewares.Auth, layananHandler.DeleteLayanan)
}