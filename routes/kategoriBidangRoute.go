package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func KategoriRoute(api fiber.Router, kategoriCollection *mongo.Collection){
	// repositories
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)

	// services
	kategoriService := services.NewKategoriBidangService(kategoriRepository)

	// controllers
	kategoriHandler := controllers.NewKategoriHandler(kategoriService)

	kategori := api.Group("/kategori")

	kategori.Post("/save", middlewares.Auth, kategoriHandler.Save)
}