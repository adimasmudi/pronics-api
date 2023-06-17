package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func BankRoute(api fiber.Router, bankCollection *mongo.Collection) {
	// repositories
	bankRepository := repositories.NewBankRepository(bankCollection)

	// services
	bankService := services.NewBankService(bankRepository)

	// controllers
	bankHandler := controllers.NewBankHandler(bankService)

	bank := api.Group("/bank")

	bank.Post("/save", middlewares.Auth, bankHandler.Save)
	bank.Get("/all", bankHandler.FindAll)
	bank.Get("/detail/:bankId", middlewares.Auth,bankHandler.FindById)
	bank.Put("/update/:bankId", middlewares.Auth, bankHandler.UpdateBank)
}