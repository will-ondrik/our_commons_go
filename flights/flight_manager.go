package flight

import (
	"etl_our_commons/dtos"
	"fmt"
)

type FlightManager struct {
	AirportService *AirportService
	FlightService *FlightService
	FlightCache dtos.FlightCache
}

func NewFlightManager() (*FlightManager, error) {
	as, err := NewAirportService()
	if err != nil {
		panic("failed to init airport service")
	}

	fs, err := NewFlightService()
	if err != nil {
		panic("failed to init airport service")
	}

	return &FlightManager{
		AirportService: as,
		FlightService: fs,
		FlightCache: make(dtos.FlightCache),
	}, nil
}


func (fm *FlightManager) GetFlightData(departureCity, destinationCity string) (*dtos.CarbonInterfaceResponse, error) {
	cities := fmt.Sprintf("%s_%s", departureCity, destinationCity)

	// Check flight cache
	if travelData := fm.FlightService.GetCache(cities); travelData != nil {
		return travelData, nil
	}

	// Retrieve airport data from cache or fetch if missing
	airports := fm.AirportService.GetCache(cities)
	if airports == nil {
		var err error
		airports, err = fm.AirportService.GetAirports(departureCity, destinationCity)
		fmt.Printf("inside nil - airports: %v", airports)
		if err != nil {
			return nil, err
		}
		fm.AirportService.SetCache(cities, *airports)
	}

	// Fetch flight estimate and update cache
	travelData, err := fm.FlightService.GetFlightEstimate(airports.DepartureAirport, airports.DestinationAirport)
	if err != nil {
		return nil, err
	}
	fm.FlightService.SetCache(cities, *travelData)

	return travelData, nil
}
