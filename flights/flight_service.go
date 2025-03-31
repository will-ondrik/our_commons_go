package flight

import (
	"bytes"
	"encoding/json"
	"etl_our_commons/dtos"
	"fmt"
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
	fmt.Println("Initializing FlightService...")
	
	err := godotenv.Load()
	if err != nil {
		fmt.Printf("Warning: Failed to load environment variables: %v\n", err)
		fmt.Println("Will attempt to use environment variables that are already set")
	}

	url := os.Getenv("CARBON_INTERFACE_URL")
	apiKey := os.Getenv("CARBON_INTERFACE_API_KEY")
	
	if url == "" {
		return nil, fmt.Errorf("CARBON_INTERFACE_URL environment variable is not set")
	}
	
	if apiKey == "" {
		return nil, fmt.Errorf("CARBON_INTERFACE_API_KEY environment variable is not set")
	}

	if !strings.Contains(url, "www.") && strings.Contains(url, "carboninterface.com") {
		url = strings.Replace(url, "carboninterface.com", "www.carboninterface.com", 1)
		fmt.Printf("Fixed API URL format: %s\n", url)
	}

	service := &FlightService{
		Cache:  make(dtos.FlightCache),
		ApiUrl: url,
		ApiKey: apiKey,
	}
	
	fmt.Println("FlightService initialized successfully")
	return service, nil
}

func (f *FlightService) GetCache(cities string) *dtos.CarbonInterfaceResponse {

	if f == nil {
		fmt.Println("Warning: FlightService is nil in GetCache")
		return nil
	}
	
	if f.Cache == nil {
		fmt.Println("Warning: Cache is nil in GetCache")
		return nil
	}
	
	// Validate input 
	if cities == "" {
		fmt.Println("Warning: cities parameter is empty in GetCache")
		return nil
	}
	
	// Check if data exists in cache
	if data, ok := f.Cache[cities]; ok {
		fmt.Printf("Found flight data in cache for %s\n", cities)
		return &data
	}
	
	fmt.Printf("No flight data found in cache for %s\n", cities)
	return nil
}

func (f *FlightService) SetCache(cities string, data dtos.CarbonInterfaceResponse) {
	if f == nil {
		fmt.Println("Warning: FlightService is nil in SetCache")
		return
	}
	
	if f.Cache == nil {
		fmt.Println("Warning: Cache is nil in SetCache, initializing new cache")
		f.Cache = make(dtos.FlightCache)
	}
	
	if cities == "" {
		fmt.Println("Warning: cities parameter is empty in SetCache")
		return
	}
	
	if data.Data == nil || data.Data.Attributes == nil {
		fmt.Println("Warning: Flight data is invalid, not caching")
		return
	}
	
	// Set data in cache
	f.Cache[cities] = data
	fmt.Printf("Cached flight data for %s\n", cities)
}

func (f *FlightService) GetFlightEstimate(departureData, destinationData dtos.AirportData) (*dtos.CarbonInterfaceResponse, error) {
	if f == nil {
		return nil, fmt.Errorf("flight service is nil")
	}
	
	if departureData.City == "" || destinationData.City == "" {
		return nil, fmt.Errorf("departure city or destination city is empty")
	}
	
	cities := fmt.Sprintf("%s_%s", departureData.City, destinationData.City)
	fmt.Printf("Getting flight estimate for %s to %s\n", departureData.City, destinationData.City)
	
	// Check cache first
	if flightData := f.GetCache(cities); flightData != nil {
		return flightData, nil
	}

	// Validate IATA codes before making the API call
	if departureData.IATA == "" || destinationData.IATA == "" {
		return nil, fmt.Errorf("missing IATA code for airports: %s, %s", departureData.City, destinationData.City)
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

	// Marshal request to JSON
	reqBody, err := json.Marshal(ciRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Add delay to avoid rate limiting
	time.Sleep(2 * time.Second)

	if f.ApiKey == "" {
		return nil, fmt.Errorf("API key is empty")
	}
	
	if f.ApiUrl == "" {
		return nil, fmt.Errorf("API URL is empty")
	}

	fmt.Printf("Making API request for flight from %s to %s\n", departureData.IATA, destinationData.IATA)
	
	// Prepare curl command
	cmd := exec.Command("curl",
		"-X", "POST",
		"-H", "Authorization: Bearer "+f.ApiKey,
		"-H", "Content-Type: application/json",
		"--data", string(reqBody),
		f.ApiUrl,
	)

	// Capture command output
	var out, stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	// Execute command
	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("cURL execution failed: %w | stderr: %s", err, stderr.String())
	}

	// Validate response
	responseBody := out.String()
	if responseBody == "" {
		return nil, fmt.Errorf("API returned an empty response. stderr: %s", stderr.String())
	}

	// Parse response JSON
	var data dtos.CarbonInterfaceResponse
	if err := json.Unmarshal([]byte(responseBody), &data); err != nil {
		return nil, fmt.Errorf("failed to parse response JSON: %w | response: %s", err, responseBody)
	}

	// Validate response data
	if data.Data == nil {
		return nil, fmt.Errorf("API response data field is nil: %s", responseBody)
	}
	
	if data.Data.ID == "" {
		return nil, fmt.Errorf("API response does not contain a valid ID: %s", responseBody)
	}
	
	if data.Data.Attributes == nil {
		return nil, fmt.Errorf("API response attributes field is nil: %s", responseBody)
	}

	// Cache the result
	f.SetCache(cities, data)
	fmt.Printf("Successfully retrieved flight estimate for %s to %s\n", departureData.City, destinationData.City)

	return &data, nil
}
