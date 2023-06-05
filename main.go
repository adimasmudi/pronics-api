package main

import (
	"fmt"
	"os"

	"pronics-api/configs"

	"github.com/gofiber/fiber"
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Println(port)

	app.Listen(":" + port)
}