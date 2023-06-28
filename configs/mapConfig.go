package configs

import (
	"log"

	"googlemaps.github.io/maps"
)

func InitMap() *maps.Client {
	c, err := maps.NewClient(maps.WithAPIKey("AIzaSyB8Xcw0-bTqcs2vXOQ5SANu65-4IR1rRFc"))

	if err != nil {
		log.Fatalf("fatal error: %s", err)
	}

	return c
}