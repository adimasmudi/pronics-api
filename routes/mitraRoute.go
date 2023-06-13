package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func MitraRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection, galeriMitraCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	galeriMitraRepository := repositories.NewGaleriRepository(galeriMitraCollection)
	
	// services
	mitraService := services.NewMitraService(userRepository, mitraRepository, galeriMitraRepository)

	// controllers
	mitraHandler := controllers.NewMitraHandler(mitraService)

	mitraRoute := api.Group("/mitra")

	mitraRoute.Get("/profile", middlewares.Auth, mitraHandler.GetProfile)
	mitraRoute.Put("/profile/update", middlewares.Auth, mitraHandler.UpdateProfile)
	mitraRoute.Put("/galeri/upload", middlewares.Auth, mitraHandler.UploadMultipleImagesToGaleri)

}