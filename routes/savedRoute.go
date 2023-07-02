package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func SavedRoute(api fiber.Router, userCollection *mongo.Collection, customerCollection *mongo.Collection, mitraCollection *mongo.Collection,wilayahCakupanCollection *mongo.Collection, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection, layananCollection *mongo.Collection, layananMitraCollection *mongo.Collection, savedCollection *mongo.Collection) {
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

	// services
	savedService := services.NewSavedService(userRepository, customerRepository, mitraRepository,wilayahCakupanRepository, bidangRepository, kategoriRepository, layananRepository, layananMitraRepository, savedRepository)

	// controllers
	savedHandler := controllers.NewSavedHandler(savedService)

	savedRoute := api.Group("/saved")

	savedRoute.Post("/add/:mitraId", middlewares.Auth, savedHandler.Save)
	savedRoute.Delete("/delete/:savedId", middlewares.Auth, savedHandler.DeleteSaved)
	savedRoute.Get("/all", middlewares.Auth, savedHandler.ShowAllSaved)
}