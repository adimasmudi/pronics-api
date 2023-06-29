package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderDetailRoute(api fiber.Router, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection) {
	// repositories

	orderRepository := repositories.NewOrderRepository(orderCollection)
	orderDetailRepository := repositories.NewOrderDetailRepository(orderDetailCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)

	// services
	orderDetailService := services.NewOrderDetailService(orderRepository, orderDetailRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository)

	// controllers
	orderDetailHandler := controllers.NewOrderDetailHandler(orderDetailService)

	orderDetail := api.Group("/orderDetail")

	orderDetail.Post("/createOrUpdate/:orderId", middlewares.Auth, orderDetailHandler.AddOrUpdateOrderDetail)
}