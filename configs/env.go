package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func EnvMONGOURI() string {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error when loading .env file", err.Error())
	}

	return os.Getenv("MONGO_URI")
}