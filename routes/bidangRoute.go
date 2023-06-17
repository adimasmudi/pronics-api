package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func BidangRoute(api fiber.Router, bidangCollection *mongo.Collection, kategoriCollection *mongo.Collection){
	// repositories
	bidangRepository := repositories.NewBidangRepository(bidangCollection)
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)

	// services
	bidangService := services.NewbidangService(bidangRepository, kategoriRepository)

	// controllers
	bidangHandler := controllers.NewbidangHandler(bidangService)

	bidang := api.Group("/bidang")

	bidang.Post("/save", middlewares.Auth, bidangHandler.Save)
	bidang.Get("/all", bidangHandler.FindAll)
	bidang.Get("/detail/:bidangId", bidangHandler.FindById)
	bidang.Put("/update/:bidangId", middlewares.Auth, bidangHandler.UpdateBidang)
	bidang.Delete("/delete/:bidangId", middlewares.Auth, bidangHandler.DeleteBidang)
}