package api

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)
type CarbonInterfaceResponse struct {
	Data CarbonInterfaceData
}

type CarbonInterfaceData struct {
	ID string `json:"id"`
	Type string `json:"type"`
	Attributes Attributes `json:"attributes"`
	EstimatedAt string `json:"estimated_at"`
	CarbonGrams int `json:"carbon_g"`
	CarbonPounds int `json:"carbon_lb"`
	CarbonKilograms int `json:"carbon_kg"`
	CarbonMetricTonnes int `json:"carbon_mt"`
	DistanceUnit string `json:"distance_unit"`
	DistanceValue float64 `json:"distance_value"`
}

type Attributes struct {
	Passengers int
	Legs Legs
}

type Legs struct {
	Flights []Flights
}

type Flights struct {
	DepartureAirport string `json:"departure_airport"`
	DestinationAirport string `json:"destination_airport"`
}


type CarbonInterfaceRequest struct {
	Type string `json:"type"`
	Passengers int `json:"passenger"`
	Legs Flights `json:"legs"`
} 


func GetFlightEstimate(departureAirportCode, destinationAirportCode string) (*CarbonInterfaceResponse, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	url := os.Getenv("CARBON_INTERFACE_URL")
	apiKey := os.Getenv("CARBON_INTERFACE_API_KEY")

	// generate req obj
	ciRequest := CarbonInterfaceRequest{
		Type: "flight",
		Passengers: 1,
		Legs: Flights{
			DepartureAirport: departureAirportCode,
			DestinationAirport: destinationAirportCode,
		},
	}


	reqBody, err := json.Marshal(ciRequest)
	if err != nil {
		log.Fatal(err)
	}


	// create new req
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type",  "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data *CarbonInterfaceResponse
	err = json.Unmarshal(body, &data)
	if err != nil {
		log.Fatal(err)
	}

	return data, nil
}