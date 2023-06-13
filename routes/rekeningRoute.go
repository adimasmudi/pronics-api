package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func RekeningRoute(api fiber.Router, rekeningCollection *mongo.Collection, bankCollection *mongo.Collection) {
	// repositories
	rekeningRepository := repositories.NewRekeningRepository(rekeningCollection)
	bankRepository := repositories.NewBankRepository(bankCollection)

	// services
	rekeningService := services.NewRekeningService(rekeningRepository, bankRepository)

	// controllers
	rekeningHandler := controllers.NewRekeningHandler(rekeningService)

	rekening := api.Group("/rekening")

	rekening.Get("/detail", middlewares.Auth, rekeningHandler.GetDetailRekening)
}