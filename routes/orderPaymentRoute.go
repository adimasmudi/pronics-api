package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderPaymentRoute(api fiber.Router, mitraCollection *mongo.Collection, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection, orderPaymentCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection) {
	// repositories

	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	orderDetailRepository := repositories.NewOrderDetailRepository(orderDetailCollection)
	orderPaymentRepository := repositories.NewOrderPaymentRepository(orderPaymentCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)

	// services
	orderPaymentService := services.NewOrderPaymentService(mitraRepository,orderRepository, orderDetailRepository,orderPaymentRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository)

	// controllers
	orderPaymentHandler := controllers.NewOrderPaymentHandler(orderPaymentService)

	orderPayment := api.Group("/orderPayment")

	orderPayment.Post("/createOrUpdate/:orderDetailId", middlewares.Auth, orderPaymentHandler.AddOrUpdateOrderPayment)
	orderPayment.Post("/confirmPayment/:orderPaymentId", middlewares.Auth, orderPaymentHandler.ConfirmPayment)
}