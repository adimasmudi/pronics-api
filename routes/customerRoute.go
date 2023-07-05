package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CustomerRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, alamatCustomerCollection *mongo.Collection, orderCollection *mongo.Collection, adminCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	alamatCustomerRepository := repositories.NewAlamatCustomerRepository(alamatCustomerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	adminRepository := repositories.NewAdminRepository(adminCollection)
	
	// services
	customerService := services.NewCustomerService(userRepository, customerRepository, alamatCustomerRepository, orderRepository)

	// controllers
	customerHandler := controllers.NewCustomerHandler(customerService)

	// auth
	customer := middlewares.CustomerAuth(customerRepository)
	admin := middlewares.AdminAuth(adminRepository)

	authCustomer := api.Group("/customer")


	authCustomer.Get("/profile", customer.AuthCustomer, customerHandler.GetProfile)
	authCustomer.Put("/profile/update", customer.AuthCustomer, customerHandler.UpdateProfil)
	authCustomer.Get("/all",admin.AuthAdmin, customerHandler.GetAllCustomer)
}