package routes

import (
	"pronics-api/controllers"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, mitraCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	customerRepository := repositories.NewcustomerRepository(customerCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)

	// services
	userService := services.NewUserService(userRepository, customerRepository, mitraRepository)

	// controllers
	userHandler := controllers.NewUserHandler(userService)

	authUser := api.Group("/auth/user")

	authUser.Post("/register", userHandler.Register)
	authUser.Post("/login", userHandler.Login)
	authUser.Post("/registerMitra", userHandler.RegisterMitra)

}