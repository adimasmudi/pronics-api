package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func KategoriRoute(api fiber.Router, kategoriCollection *mongo.Collection, bidangCollection *mongo.Collection){
	// repositories
	kategoriRepository := repositories.NewKategoriRepository(kategoriCollection)
	bidangRepository := repositories.NewBidangRepository(bidangCollection)

	// services
	kategoriService := services.NewKategoriBidangService(kategoriRepository, bidangRepository)

	// controllers
	kategoriHandler := controllers.NewKategoriHandler(kategoriService)

	kategori := api.Group("/kategori")

	kategori.Post("/save", middlewares.Auth, kategoriHandler.Save)
	kategori.Get("/all", kategoriHandler.FindAll)
	kategori.Get("/all/bidang", kategoriHandler.GetKategoriWithBidang)
	kategori.Get("/:kategoriId", kategoriHandler.GetKategoriById)
}