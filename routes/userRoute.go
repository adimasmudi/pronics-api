package routes

import (
	"pronics-api/controllers"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)

	// services
	userService := services.NewUserService(userRepository)

	// controllers
	userHandler := controllers.NewUserHandler(userService)

	authUser := api.Group("/auth/user")

	authUser.Post("/register", userHandler.Register)
	authUser.Post("/login", userHandler.Login)

}