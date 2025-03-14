package airports

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type AirportsList map[string]AirportData


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

type Trip struct {
	DepartureAirport AirportData
	DestinationAirport AirportData
}

type AirportCache map[string]AirportData

type AirportDetails struct {
	DetailsMap AirportsList
	Cache AirportCache
}


func NewAirportDetails() (*AirportDetails, error) {
	airportsMap, err := ParseAirportsFile()
	if err != nil {
		return nil, err
	}

	return &AirportDetails{
		DetailsMap: airportsMap,
		Cache: make(map[string]AirportData),
	}, nil
}


func ParseAirportsFile() (AirportsList, error) {
	file, err := os.Open("./airports/airports.json")
	if err != nil {
		log.Fatal("failed to open airports.json")
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("failed to read file")
	}


	var airports AirportsList
	err = json.Unmarshal(bytes, &airports)
	if err != nil {
		log.Fatal("failed to unmarshal airport list")
	}
	
	return airports, nil
}



// TODO: Remove hardcoded country code
func (a *AirportDetails) GetAirportDetails(departureCity, destinationCity string) (Trip, error) {
	airportsList, err := ParseAirportsFile()
	if err != nil {
		return Trip{}, err
	}

	trip := Trip{}
	for _, airport := range airportsList {
		if airport.City == departureCity && airport.Country == "CA" {
			trip.DepartureAirport = airport
		}

		if airport.City == destinationCity && airport.Country == "CA" {
			trip.DestinationAirport = airport
		}

		if trip.DepartureAirport != (AirportData{}) && trip.DestinationAirport != (AirportData{}) {

		}
	}
	return trip, nil
}

// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportDetails) GetCache(cities string) *AirportData {
	if data, ok := a.Cache[cities]; ok {
		return &data
	} else {
		return nil
	}
}

// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportDetails) SetCache (cities string, airportData AirportData) {
	a.Cache[cities] = airportData
}



