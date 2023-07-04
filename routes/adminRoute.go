package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AdminRoute(api fiber.Router, adminCollection *mongo.Collection, mitraCollection *mongo.Collection, customerCollection *mongo.Collection, orderCollection *mongo.Collection) {
	// repositories
	adminRepository := repositories.NewAdminRepository(adminCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)

	// services
	adminService := services.NewAdminService(adminRepository, mitraRepository, customerRepository, orderRepository)

	// controllers
	adminHandler := controllers.NewAdminHandler(adminService)

	authAdmin := api.Group("/auth/admin")

	authAdmin.Post("/register", adminHandler.Register)
	authAdmin.Post("/login", adminHandler.Login)
	authAdmin.Get("/profile", middlewares.Auth,adminHandler.GetProfile)
	authAdmin.Get("/dashboardSummary", middlewares.Auth, adminHandler.GetDashboardSummary)
}