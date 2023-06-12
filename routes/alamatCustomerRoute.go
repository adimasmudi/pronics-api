package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func AlamatCustomerRoute(api fiber.Router, alamatCustomerCollection *mongo.Collection, customerCollection *mongo.Collection) {
	// repositories
	alamatCustomerRepository := repositories.NewAlamatCustomerRepository(alamatCustomerCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)

	// services
	alamatCustomerService := services.NewAlamatCustomerService(alamatCustomerRepository, customerRepository)

	// controllers
	alamatCustomerHandler := controllers.NewAlamatCustomerHandler(alamatCustomerService)

	alamat := api.Group("/alamat")

	alamat.Post("/save", middlewares.Auth, alamatCustomerHandler.Save)
}