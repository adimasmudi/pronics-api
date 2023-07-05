package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AlamatCustomerRoute(api fiber.Router, alamatCustomerCollection *mongo.Collection, customerCollection *mongo.Collection, userCollection *mongo.Collection) {
	// repositories
	alamatCustomerRepository := repositories.NewAlamatCustomerRepository(alamatCustomerCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	userRepository := repositories.NewUserRepository(userCollection)

	// services
	alamatCustomerService := services.NewAlamatCustomerService(alamatCustomerRepository, customerRepository, userRepository)

	// controllers
	alamatCustomerHandler := controllers.NewAlamatCustomerHandler(alamatCustomerService)

	// auth
	customer := middlewares.CustomerAuth(customerRepository)

	alamat := api.Group("/alamat")

	alamat.Post("/save", customer.AuthCustomer, alamatCustomerHandler.Save)
	alamat.Get("/all", customer.AuthCustomer, alamatCustomerHandler.GetAllAlamatCustomer)
	alamat.Put("/setAsUtama/:alamatId", customer.AuthCustomer, alamatCustomerHandler.SetAsAlamatUtama)
}