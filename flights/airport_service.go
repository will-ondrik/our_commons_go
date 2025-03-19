package flight

import (
	"encoding/json"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"io"
	"log"
	"os"
	"slices"
)


type AirportService struct {
	AirportsMap dtos.AirportsList
	Cache dtos.AirportCache
}


func NewAirportService() (*AirportService, error) {
	airportsMap, err := ParseAirportsFile()
	if err != nil {
		return nil, fmt.Errorf("failed to parse airports.json: %v", err)
	}

	return &AirportService{
		AirportsMap: airportsMap,
		Cache: make(dtos.AirportCache),
	}, nil
}


func ParseAirportsFile() (dtos.AirportsList, error) {
	file, err := os.Open("./flights/airports.json")
	if err != nil {
		log.Fatal("failed to open airports.json")
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		log.Fatal("failed to read file")
	}


	var airports dtos.AirportsList
	err = json.Unmarshal(bytes, &airports)
	if err != nil {
		log.Fatal("failed to unmarshal airport list")
	}
	
	return airports, nil
}



// TODO: Remove hardcoded country code
func (a *AirportService) GetAirports(departureCity, destinationCity string) (*dtos.Trip, error) {

	trip := dtos.Trip{}
	for _, airport := range a.AirportsMap {
		if airport.City == departureCity && airport.Country == "CA" {
			trip.DepartureAirport = airport
		}

		if airport.City == destinationCity && airport.Country == "CA" {
			trip.DestinationAirport = airport
		}

		if trip.DepartureAirport != (dtos.AirportData{}) && trip.DestinationAirport != (dtos.AirportData{}) {
			fmt.Println("retruning here", trip)
			return &trip, nil
		}
	}
	return &trip, nil
}

// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportService) GetCache(cities string) *dtos.Trip {
	if data, ok := a.Cache[cities]; ok {
		return &data
	} 

	return nil
}

// Expects DepartureCityCountry_DestinationCityCountry
// Example: VancouverCanada_OttawaCanada
func (a *AirportService) SetCache (cities string, airportData dtos.Trip) {
	a.Cache[cities] = airportData
}



func (a *AirportService) IsFlight(travelPurpose string) bool {
	// to attend a national caucus meeting
	// to attend a regional or provincial caucus meeting
	// Attending event with Member (type: Employee)
	// // Need to compare cost to ensure its a flight
		// There may be multiple cities in close proximity
	// unite the family with the Member
	// travel to/from constituency and Ottawa

	// Potentials:
		// to attend training
		// to attend meetings with stakeholders about business of the House
		// to support a parliamentary association
		// to attend language training
		//

	return slices.Contains(constants.FLIGHT_KEYWORDS, travelPurpose)
}