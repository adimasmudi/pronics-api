package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AdminRoute(api fiber.Router, adminCollection *mongo.Collection) {
	// repositories
	adminRepository := repositories.NewAdminRepository(adminCollection)

	// services
	adminService := services.NewAdminService(adminRepository)

	// controllers
	adminHandler := controllers.NewAdminHandler(adminService)

	authAdmin := api.Group("/auth/admin")

	authAdmin.Post("/register", adminHandler.Register)
	authAdmin.Post("/login", adminHandler.Login)
	authAdmin.Get("/profile", middlewares.Auth,adminHandler.GetProfile)
}