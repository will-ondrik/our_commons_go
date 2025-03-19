package flight

import (
	"bytes"
	"encoding/json"
	"etl_our_commons/dtos"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type FlightService struct {
	Cache  dtos.FlightCache
	ApiUrl string
	ApiKey string
}

func NewFlightService() (*FlightService, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Failed to load environment variables")
	}

	url := os.Getenv("CARBON_INTERFACE_URL")
	apiKey := os.Getenv("CARBON_INTERFACE_API_KEY")

	if !strings.Contains(url, "www.") {
		url = strings.Replace(url, "carboninterface.com", "www.carboninterface.com", 1)
	}

	return &FlightService{
		Cache:  make(dtos.FlightCache),
		ApiUrl: url,
		ApiKey: apiKey,
	}, nil
}

func (f *FlightService) GetCache(cities string) *dtos.CarbonInterfaceResponse {
	if data, ok := f.Cache[cities]; ok {
		return &data
	}
	return nil
}

func (f *FlightService) SetCache(cities string, data dtos.CarbonInterfaceResponse) {
	f.Cache[cities] = data
}

func (f *FlightService) GetFlightEstimate(departureData, destinationData dtos.AirportData) (*dtos.CarbonInterfaceResponse, error) {
	cities := fmt.Sprintf("%s_%s", departureData.City, destinationData.City)
	if flightData := f.GetCache(cities); flightData != nil {
		return flightData, nil
	}

	ciRequest := dtos.CarbonInterfaceRequest{
		Type:       "flight",
		Passengers: 1,
		Legs: []dtos.Flights{
			{
				DepartureAirport:   departureData.IATA,
				DestinationAirport: destinationData.IATA,
			},
		},
	}

	reqBody, err := json.Marshal(ciRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	time.Sleep(2 * time.Second)

	if f.ApiKey == "" || f.ApiUrl == "" {
		return nil, fmt.Errorf("missing API key or URL")
	}

	cmd := exec.Command("curl",
		"-X", "POST",
		"-H", "Authorization: Bearer "+f.ApiKey,
		"-H", "Content-Type: application/json",
		"--data", string(reqBody),
		f.ApiUrl,
	)

	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("cURL execution failed: %w | stderr: %s", err, stderr.String())
	}

	responseBody := out.String()
	if responseBody == "" {
		return nil, fmt.Errorf("API returned an empty response. stderr: %s", stderr.String())
	}

	var data dtos.CarbonInterfaceResponse
	if err := json.Unmarshal([]byte(responseBody), &data); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w", err)
	}

	if data.Data.ID == "" {
		return nil, fmt.Errorf("API response does not contain a valid ID")
	}

	f.SetCache(cities, data)

	return &data, nil
}
