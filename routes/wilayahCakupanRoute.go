package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func WilayahCakupanRoute(api fiber.Router, wilayahCakupanCollection *mongo.Collection){
	// repositories
	wilayahCakupanRepository := repositories.NewWilayahRepository(wilayahCakupanCollection)

	// services
	wilayahCakupanService := services.NewWilayahCakupanService(wilayahCakupanRepository)

	// controllers
	wilayahCakupanHandler := controllers.NewwilayahCakupanHandler(wilayahCakupanService)

	WilayahCakupan := api.Group("/wilayahCakupan")

	WilayahCakupan.Post("/save", middlewares.Auth, wilayahCakupanHandler.Save)
	WilayahCakupan.Get("/all", wilayahCakupanHandler.FindAll)
}