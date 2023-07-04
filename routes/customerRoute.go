package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func CustomerRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, alamatCustomerCollection *mongo.Collection, orderCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	alamatCustomerRepository := repositories.NewAlamatCustomerRepository(alamatCustomerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	
	// services
	customerService := services.NewCustomerService(userRepository, customerRepository, alamatCustomerRepository, orderRepository)

	// controllers
	customerHandler := controllers.NewCustomerHandler(customerService)

	authCustomer := api.Group("/customer")

	authCustomer.Get("/profile", middlewares.Auth, customerHandler.GetProfile)
	authCustomer.Put("/profile/update", middlewares.Auth, customerHandler.UpdateProfil)
	authCustomer.Get("/all", middlewares.Auth, customerHandler.GetAllCustomer)
}