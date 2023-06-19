package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func LayananMitraRoute(api fiber.Router, layananMitraCollection *mongo.Collection, bidangCollection *mongo.Collection, mitraCollection *mongo.Collection) {
	// repositories
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)

	// services
	layananMitraService := services.NewLayananMitraService(layananMitraRepository, bidangRepository, mitraRepository)

	// controllers
	layananMitraHandler := controllers.NewLayananMitraHandler(layananMitraService)

	layananMitra := api.Group("/layananMitra")

	layananMitra.Post("/save", middlewares.Auth, layananMitraHandler.Save)
	layananMitra.Delete("/delete/:layananMitraId", middlewares.Auth, layananMitraHandler.Delete)
}
