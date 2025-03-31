package flight

import (
	"encoding/json"
	"etl_our_commons/dtos"
	"fmt"
	"io"
	"math"
	"os"
	"strings"
)


type AirportService struct {
	AirportsMap dtos.AirportsList
	Cache dtos.AirportCache
}


func NewAirportService() (*AirportService, error) {
	fmt.Println("Initializing AirportService...")
	
	airportsMap, err := ParseAirportsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to parse airports.json: %w", err)
	}
	
	// Validate airports map
	if airportsMap == nil {
		return nil, fmt.Errorf("airports map is nil after parsing")
	}
	
	if len(airportsMap) == 0 {
		fmt.Println("Warning: Airports map is empty")
	} else {
		fmt.Printf("Successfully loaded %d airports\n", len(airportsMap))
	}

	service := &AirportService{
		AirportsMap: airportsMap,
		Cache: make(dtos.AirportCache),
	}
	
	fmt.Println("AirportService initialized successfully")
	return service, nil
}

func ParseAirportsFile() (dtos.AirportsList, error) {
	fmt.Println("Parsing airports file...")
	
	file, err := os.Open("./flights/airports.json")
	if err != nil {
		return nil, fmt.Errorf("failed to open airports.json: %w", err)
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read airports.json: %w", err)
	}
	
	// Check if file is empty
	if len(bytes) == 0 {
		return nil, fmt.Errorf("airports.json is empty")
	}

	var airports dtos.AirportsList
	err = json.Unmarshal(bytes, &airports)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal airport list: %w", err)
	}
	
	if airports == nil {
		return nil, fmt.Errorf("parsed airports map is nil")
	}
	
	if len(airports) == 0 {
		fmt.Println("Warning: No airports found in airports.json")
	} else {
		fmt.Printf("Successfully parsed %d airports from file\n", len(airports))
	}
	
	return airports, nil
}


func (a *AirportService) GetAirports(departureCity, destinationCity string) (*dtos.Trip, error) {
	if a == nil {
		return nil, fmt.Errorf("airport service is nil")
	}
	if a.AirportsMap == nil {
		return nil, fmt.Errorf("airports map is nil")
	}
	if departureCity == "" || destinationCity == "" {
		return nil, fmt.Errorf("departure city or destination city is empty")
	}

	fmt.Printf("Looking for airports in %s and %s\n", departureCity, destinationCity)

	trip := dtos.Trip{}
	foundDeparture := false
	foundDestination := false

	for _, airport := range a.AirportsMap {

		// Handle departure airport
		if !foundDeparture && (airport.City == departureCity || strings.Contains(airport.Name, departureCity)) {
			if airport.IATA != "" {
				trip.DepartureAirport = airport
				foundDeparture = true
				fmt.Printf("Found departure airport for %s: %s (%s)\n", departureCity, airport.Name, airport.IATA)
			} else {
				fmt.Printf("Warning: No IATA code for departure airport in %s, searching nearby...\n", departureCity)
				nearest := a.FindNearestAirport(airport, a.AirportsMap)
				if nearest.IATA != "" {
					trip.DepartureAirport = nearest
					foundDeparture = true
					fmt.Printf("Using nearest airport for %s: %s (%s)\n", departureCity, nearest.Name, nearest.IATA)
				}
			}
		}

		// Handle destination airport
		if !foundDestination && airport.City == destinationCity && airport.Country == "CA" {
			if airport.IATA != "" {
				trip.DestinationAirport = airport
				foundDestination = true
				fmt.Printf("Found destination airport for %s: %s (%s)\n", destinationCity, airport.Name, airport.IATA)
			} else {
				fmt.Printf("Warning: No IATA code for destination airport in %s, searching nearby...\n", destinationCity)
				nearest := a.FindNearestAirport(airport, a.AirportsMap)
				if nearest.IATA != "" {
					trip.DestinationAirport = nearest
					foundDestination = true
					fmt.Printf("Using nearest airport for %s: %s (%s)\n", destinationCity, nearest.Name, nearest.IATA)
				}
			}
		}

		if foundDeparture && foundDestination {
			return &trip, nil
		}
	}

	if !foundDeparture && !foundDestination {
		return nil, fmt.Errorf("could not find airports for both %s and %s", departureCity, destinationCity)
	}
	if !foundDeparture {
		return nil, fmt.Errorf("could not find valid IATA airport for %s", departureCity)
	}
	if !foundDestination {
		return nil, fmt.Errorf("could not find valid IATA airport for %s", destinationCity)
	}

	return &trip, nil
}


// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportService) GetCache(cities string) *dtos.Trip {
	// Check if AirportService is properly initialized
	if a == nil {
		fmt.Println("Warning: AirportService is nil in GetCache")
		return nil
	}
	
	// Check if Cache is initialized
	if a.Cache == nil {
		fmt.Println("Warning: Cache is nil in GetCache")
		return nil
	}
	
	// Validate input parameter
	if cities == "" {
		fmt.Println("Warning: cities parameter is empty in GetCache")
		return nil
	}
	
	// Check if data exists in cache
	if data, ok := a.Cache[cities]; ok {
		fmt.Printf("Found airport data in cache for %s\n", cities)
		return &data
	} 

	fmt.Printf("No airport data found in cache for %s\n", cities)
	return nil
}

// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportService) SetCache(cities string, airportData dtos.Trip) {
	// Check if AirportService is properly initialized
	if a == nil {
		fmt.Println("Warning: AirportService is nil in SetCache")
		return
	}
	
	// Check if Cache is initialized
	if a.Cache == nil {
		fmt.Println("Warning: Cache is nil in SetCache, initializing new cache")
		a.Cache = make(dtos.AirportCache)
	}
	
	// Validate input parameters
	if cities == "" {
		fmt.Println("Warning: cities parameter is empty in SetCache")
		return
	}
	
	// Validate airport data has valid IATA codes
	if airportData.DepartureAirport.IATA == "" || airportData.DestinationAirport.IATA == "" {
		fmt.Println("Warning: Airport data has missing IATA codes, not caching")
		return
	}
	
	// Set data in cache
	a.Cache[cities] = airportData
	fmt.Printf("Cached airport data for %s\n", cities)
}


func haversine(lat1, lon1, lat2, lon2 float64) float64 {
	const RAD = 6371 // Earth radius in km

	dLat := (lat2 - lat1) * (math.Pi / 180)
	dLon := (lon2 - lon1) * (math.Pi / 180)

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*(math.Pi/180))*math.Cos(lat2*(math.Pi/180))*
			math.Sin(dLon/2)*math.Sin(dLon/2)

	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return RAD * c
}

func (a *AirportService) FindNearestAirport(airport dtos.AirportData, airportsList dtos.AirportsList) dtos.AirportData {
	
	fmt.Printf("Finding nearest airport to %s (%s)\n", airport.City, airport.ICOA)
	
	var nearest dtos.AirportData
	minDistance := math.MaxFloat64
	foundValidAirport := false

	// Search for the nearest airport with a valid IATA code
	for _, ap := range airportsList {

		// Skip airports without IATA codes or in different states/timezones
		if ap.IATA == "" || ap.State != airport.State || ap.Timezone != airport.Timezone {
			continue
		}

		// Skip comparing airport to itself
		if ap.IATA == airport.IATA && ap.Name == airport.Name && ap.City == airport.City {
			continue
		}
		
		// Calculate distance between airports
		distance := haversine(ap.Latitude, ap.Longitude, airport.Latitude, airport.Longitude)
		
		// Update nearest airport if this one is closer
		if distance < minDistance {
			minDistance = distance
			nearest = ap // Use the airport from the list, not the original airport
			foundValidAirport = true
			fmt.Printf("Found closer airport: %s (%s) at distance %.2f km\n", ap.Name, ap.IATA, distance)
		}
	}
	
	// If we found a valid airport, return it
	if foundValidAirport {
		fmt.Printf("Closest airport to %s is %s (%s) at distance %.2f km\n", 
			airport.City, nearest.Name, nearest.IATA, minDistance)
		return nearest
	}
	
	// If no valid airport was found, return the original airport as fallback
	fmt.Printf("Warning: No valid airport with IATA code found near %s, using original airport\n", airport.City)
	return airport
}
