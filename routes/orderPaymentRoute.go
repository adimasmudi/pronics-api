package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderPaymentRoute(api fiber.Router,userCollection *mongo.Collection, mitraCollection *mongo.Collection,customerCollection *mongo.Collection, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection, orderPaymentCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection) {
	// repositories

	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	orderDetailRepository := repositories.NewOrderDetailRepository(orderDetailCollection)
	orderPaymentRepository := repositories.NewOrderPaymentRepository(orderPaymentCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)

	// services
	orderPaymentService := services.NewOrderPaymentService(userRepository,mitraRepository,customerRepository,orderRepository, orderDetailRepository,orderPaymentRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository)

	// controllers
	orderPaymentHandler := controllers.NewOrderPaymentHandler(orderPaymentService)

	// auth
	customer := middlewares.CustomerAuth(customerRepository)

	orderPayment := api.Group("/orderPayment")

	orderPayment.Post("/createOrUpdate/:orderDetailId", customer.AuthCustomer, orderPaymentHandler.AddOrUpdateOrderPayment)
	orderPayment.Post("/confirmPayment/:orderPaymentId", customer.AuthCustomer, orderPaymentHandler.ConfirmPayment)
}