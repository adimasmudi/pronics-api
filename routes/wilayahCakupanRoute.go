package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func WilayahCakupanRoute(api fiber.Router, wilayahCakupanCollection *mongo.Collection, adminCollection *mongo.Collection){
	// repositories
	wilayahCakupanRepository := repositories.NewWilayahRepository(wilayahCakupanCollection)
	adminRepository := repositories.NewAdminRepository(adminCollection)

	// services
	wilayahCakupanService := services.NewWilayahCakupanService(wilayahCakupanRepository)

	// controllers
	wilayahCakupanHandler := controllers.NewwilayahCakupanHandler(wilayahCakupanService)

	// auth
	adminAuth := middlewares.AdminAuth(adminRepository)

	WilayahCakupan := api.Group("/wilayahCakupan")

	WilayahCakupan.Post("/save", adminAuth.AuthAdmin, wilayahCakupanHandler.Save)
	WilayahCakupan.Get("/all", wilayahCakupanHandler.FindAll)
}