package main

import (
	"fmt"
	"log"
	"os"

	"pronics-api/configs"
	"pronics-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
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

	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file", err.Error())
	}

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
	var ktpMitraCollection *mongo.Collection = configs.GetCollection(configs.DB, "ktpMitras")
	var alamatCustomerCollection *mongo.Collection = configs.GetCollection(configs.DB, "alamatCustomer")
	var bankCollection *mongo.Collection = configs.GetCollection(configs.DB, "banks")
	var galeriMitraCollection *mongo.Collection = configs.GetCollection(configs.DB, "galeriMitra")

	api := app.Group("/api/v1")
	
	// routes
	routes.AdminRoute(api, adminCollection)
	routes.UserRoute(api, userCollection, customerCollection, mitraCollection, rekeningCollection, ktpMitraCollection)
	routes.CustomerRoute(api, userCollection, customerCollection, alamatCustomerCollection)
	routes.MitraRoute(api, userCollection, mitraCollection, galeriMitraCollection, wilayahCakupanCollection, bidangCollection)
	routes.KategoriRoute(api, kategoriCollection, bidangCollection)
	routes.WilayahCakupanRoute(api, wilayahCakupanCollection)
	routes.BidangRoute(api, bidangCollection, kategoriCollection)
	routes.AlamatCustomerRoute(api, alamatCustomerCollection, customerCollection, userCollection)
	routes.BankRoute(api, bankCollection)
	routes.RekeningRoute(api, rekeningCollection, bankCollection)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	fmt.Println("listen to port :",port)


	app.Get("/",func(c *fiber.Ctx) error {
		c.Set(fiber.HeaderContentType, fiber.MIMETextHTML)
		return c.SendString(`<html>
		<body>
			<a href="/api/v1/auth/user/login/google">Login dengan Google</a>
		</body>
		</html>`)
	})

	app.Listen(":" + port)
}