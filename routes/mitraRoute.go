package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func MitraRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	
	// services
	mitraService := services.NewMitraService(userRepository, mitraRepository)

	// controllers
	mitraHandler := controllers.NewMitraHandler(mitraService)

	mitraRoute := api.Group("/mitra")

	mitraRoute.Get("/profile", middlewares.Auth, mitraHandler.GetProfile)
	mitraRoute.Put("/profile/update", middlewares.Auth, mitraHandler.UpdateProfile)

}