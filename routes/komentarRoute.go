package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func KomentarRoute(api fiber.Router, userCollection *mongo.Collection, mitraCollection *mongo.Collection, customerCollection *mongo.Collection, orderCollection *mongo.Collection, orderDetailCollection *mongo.Collection, komentarCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection){
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	orderRepository := repositories.NewOrderRepository(orderCollection)
	orderDetailRepository := repositories.NewOrderDetailRepository(orderDetailCollection)
	komentarRepository := repositories.NewKomentarRepository(komentarCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)

	// services
	komentarService := services.NewKomentarService(userRepository, mitraRepository, customerRepository, orderRepository, orderDetailRepository, komentarRepository, layananRepository, layananMitraRepository)

	// controllers
	komentarHandler := controllers.NewKomentarHandler(komentarService)

	// auth
	customer := middlewares.CustomerAuth(customerRepository)
	authAll := middlewares.AuthAll()

	komentar := api.Group("/komentar")

	komentar.Post("/add/:orderId", customer.AuthCustomer, komentarHandler.AddKomentar)
	komentar.Get("/see/:orderId", authAll.AuthAll, komentarHandler.KomentarDetail)
	komentar.Patch("/update/:komentarId", customer.AuthCustomer, komentarHandler.UpdateKomentar)
	komentar.Patch("/response/:komentarId", customer.AuthCustomer, komentarHandler.ResponseKomentar)
	komentar.Delete("/delete/:komentarId", customer.AuthCustomer, komentarHandler.DeleteKomentar)
}