package main

import (
	"fmt"
	"os"

	"pronics-api/configs"
	"pronics-api/routes"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/mongo"
	// "github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin,Authorization",
	// 	AllowOrigins:     "*",
	// 	AllowCredentials: true,
	// 	AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	// }))

	// run database
	configs.ConnectDB()

	// collections
	var adminCollection *mongo.Collection = configs.GetCollection(configs.DB, "admins")
	var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
	var customerCollection *mongo.Collection = configs.GetCollection(configs.DB, "customers")
	var mitraCollection *mongo.Collection = configs.GetCollection(configs.DB, "mitras")
	
	
	// routes
	routes.AdminRoute(app, adminCollection)
	routes.UserRoute(app, userCollection, customerCollection, mitraCollection)
	

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port)

	app.Listen(":" + port)
}