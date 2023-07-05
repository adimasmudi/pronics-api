package routes

import (
	"pronics-api/controllers"
	"pronics-api/middlewares"
	"pronics-api/repositories"
	"pronics-api/services"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
)

func BankRoute(api fiber.Router, bankCollection *mongo.Collection, adminCollection *mongo.Collection) {
	// repositories
	bankRepository := repositories.NewBankRepository(bankCollection)
	adminRepository := repositories.NewAdminRepository(adminCollection)

	// services
	bankService := services.NewBankService(bankRepository)

	// controllers
	bankHandler := controllers.NewBankHandler(bankService)

	// auth
	adminAuth := middlewares.AdminAuth(adminRepository)

	bank := api.Group("/bank")

	bank.Post("/save", adminAuth.AuthAdmin, bankHandler.Save)
	bank.Get("/all", bankHandler.FindAll)
	bank.Get("/detail/:bankId", adminAuth.AuthAdmin,bankHandler.FindById)
	bank.Put("/update/:bankId", adminAuth.AuthAdmin, bankHandler.UpdateBank)
	bank.Delete("/delete/:bankId", adminAuth.AuthAdmin, bankHandler.DeleteBank)
}