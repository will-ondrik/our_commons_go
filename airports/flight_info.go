package airpoirts

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type AirportsList struct {
	Airports map[string]AirportData `json:"-"`
}

type AirportData struct {
	ICOA string `json:"icao"`
	IATA string `json:"iata"`
	Name string `json:"name"`
	City string `json:"city"`
	State string `json:"state"`
	Country string `json:"country"`
	Elevation int `json:"elevation"`
	Latitude float64 `json:"lat"`
	Longitude float64 `json:"lon"`
	Timezone string `json:"tz"`
}


func parseAirportsFile() (*AirportsList, error) {
	file, err := os.Open("./airports.json")
	if err != nil {
		log.Fatal("failed to open airports.json")
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("failed to read file")
	}

	var airports *AirportsList
	err = json.Unmarshal(bytes, &airports)
	if err != nil {
		log.Fatal("failed to unmarshal airport list")
	}

	return airports, nil
}

