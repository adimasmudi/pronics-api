package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection, customerCollection *mongo.Collection, orderCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)

	// services
	orderService := services.NewOrderService(userRepository, mitraRepository, customerRepository, orderRepository)

	// controllers
	orderHandler := controllers.NewOrderHandler(orderService)

	order := api.Group("/order")

	order.Post("/createTemporary/:mitraId", middlewares.Auth,orderHandler.CreateTemporaryOrder)
}