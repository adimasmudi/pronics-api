package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SavedRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, mitraCollection *mongo.Collection,wilayahCakupanCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection, savedCollection *mongo.Collection, komentarCollection *mongo.Collection) {
	// repositories
	userRepository := repositories.NewUserRepository(userCollection)
	customerRepository := repositories.NewCustomerRepository(customerCollection)
	mitraRepository := repositories.NewMitraRepository(mitraCollection)
	wilayahCakupanRepository := repositories.NewWilayahRepository(wilayahCakupanCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	layananRepository := repositories.NewLayananRepository(layananCollection)
	layananMitraRepository := repositories.NewLayananMitraRepository(layananMitraCollection)
	savedRepository := repositories.NewSavedRepository(savedCollection)
	komentarRepository := repositories.NewKomentarRepository(komentarCollection)

	// services
	savedService := services.NewSavedService(userRepository, customerRepository, mitraRepository,wilayahCakupanRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, savedRepository, komentarRepository)

	// controllers
	savedHandler := controllers.NewSavedHandler(savedService)

	// auth
	customer := middlewares.CustomerAuth(customerRepository)

	savedRoute := api.Group("/saved")

	savedRoute.Post("/add/:mitraId", customer.AuthCustomer, savedHandler.Save)
	savedRoute.Delete("/delete/:savedId", customer.AuthCustomer, savedHandler.DeleteSaved)
	savedRoute.Get("/all", customer.AuthCustomer, savedHandler.ShowAllSaved)
}