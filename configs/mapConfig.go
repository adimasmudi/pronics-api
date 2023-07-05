package configs

import (
	"log"
	"os"

	"googlemaps.github.io/maps"
)

func InitMap() *maps.Client {
	c, err := maps.NewClient(maps.WithAPIKey(os.Getenv("MAPS_API_KEY")))

	if err != nil {
		log.Fatalf("fatal error : %s", err)
	}

	return c
}