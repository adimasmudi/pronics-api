package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func OrderRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection, customerCollection *mongo.Collection, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection,orderPaymentCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection) {
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
	orderService := services.NewOrderService(userRepository, mitraRepository, customerRepository, orderRepository, orderDetailRepository,orderPaymentRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository)

	// controllers
	orderHandler := controllers.NewOrderHandler(orderService)

	order := api.Group("/order")

	order.Post("/createTemporary/:mitraId", middlewares.Auth,orderHandler.CreateTemporaryOrder)
	order.Get("/all", middlewares.Auth, orderHandler.FindAll)
	order.Get("/detail/:orderId",middlewares.Auth, orderHandler.GetOrderDetail)
	order.Patch("/updateStatus/:orderId", middlewares.Auth, orderHandler.UpdateStatus)
	order.Get("/getByMitra", middlewares.Auth, orderHandler.FindAllOrderMitra)
	order.Get("/maps/direction/:orderId", middlewares.Auth, orderHandler.GetDirection)
	order.Get("/history/all", middlewares.Auth, orderHandler.FindAllOrderHistory)
}