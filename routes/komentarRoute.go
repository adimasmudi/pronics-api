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

	komentar := api.Group("/komentar")

	komentar.Post("/add/:orderId", middlewares.Auth, komentarHandler.AddKomentar)
	komentar.Get("/see/:orderId", middlewares.Auth, komentarHandler.KomentarDetail)
	komentar.Patch("/update/:komentarId", middlewares.Auth, komentarHandler.UpdateKomentar)
	komentar.Patch("/response/:komentarId", middlewares.Auth, komentarHandler.ResponseKomentar)
}