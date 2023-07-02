package helper

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func DistanceCalculation(origin string, destination string) (float64, error) {
	url := fmt.Sprintf("https://maps.googleapis.com/maps/api/distancematrix/json?origins=%s&destinations=%s&units=imperial&key=%s", origin, destination, os.Getenv("MAPS_API_KEY"))
	url = strings.Replace(url, " ", "%20", -1)
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return 0, err
	}
	res, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)

	if err != nil {
		return 0, err
	}

	var distanceMatrix DistanceMatrixResult

	err = json.Unmarshal(body, &distanceMatrix)
	if err != nil {
		return 0, err
	}

	return float64(distanceMatrix.Rows[0].Elements[0].Distance.Value) / 1000, nil
}