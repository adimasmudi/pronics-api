package main

import (
	"os"

	"pronics-api/configs"
	"pronics-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	app := fiber.New()

	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))

	// run database
	configs.ConnectDB()

	// collections
	var adminCollection *mongo.Collection = configs.GetCollection(configs.DB, "admins")
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	var customerCollection *mongo.Collection = configs.GetCollection(configs.DB, "customers")
	var mitraCollection *mongo.Collection = configs.GetCollection(configs.DB, "mitras")
	var rekeningCollection *mongo.Collection = configs.GetCollection(configs.DB, "rekenings")
	var kategoriCollection *mongo.Collection = configs.GetCollection(configs.DB, "categories")
	var wilayahCakupanCollection *mongo.Collection = configs.GetCollection(configs.DB, "wilayahCakupans")
	var bidangCollection *mongo.Collection = configs.GetCollection(configs.DB, "bidangs")

	api := app.Group("/api/v1")
	
	// routes
	routes.AdminRoute(api, adminCollection)
	routes.UserRoute(api, userCollection, customerCollection, mitraCollection, rekeningCollection)
	routes.CustomerRoute(api, userCollection, customerCollection)
	routes.MitraRoute(api, userCollection, mitraCollection)
	routes.KategoriRoute(api, kategoriCollection, bidangCollection)
	routes.WilayahCakupanRoute(api, wilayahCakupanCollection)
	routes.BidangRoute(api, bidangCollection, kategoriCollection)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	app.Listen(":" + port)
}