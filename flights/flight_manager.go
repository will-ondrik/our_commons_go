package flight

import (
	"encoding/json"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"etl_our_commons/extract"
	"fmt"
	"os"
)

type FlightManager struct {
	AirportService *AirportService
	FlightService *FlightService
	FlightCache dtos.FlightCache
}

func NewFlightManager() (*FlightManager, error) {
	as, err := NewAirportService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize airport service: %w", err)
	}
	
	if as == nil {
		return nil, fmt.Errorf("airport service is nil after initialization")
	}

	fs, err := NewFlightService()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize flight service: %w", err)
	}
	
	if fs == nil {
		return nil, fmt.Errorf("flight service is nil after initialization")
	}

	fm := &FlightManager{
		AirportService: as,
		FlightService: fs,
		FlightCache: make(dtos.FlightCache),
	}
	
	fmt.Println("Flight manager initialized successfully")
	return fm, nil
}


func (fm *FlightManager) GetFlightData(departureCity, destinationCity string) (*dtos.CarbonInterfaceResponse, error) {
	if fm == nil {
		return nil, fmt.Errorf("flight manager is nil")
	}
	
	if fm.AirportService == nil {
		return nil, fmt.Errorf("airport service is nil")
	}
	
	if fm.FlightService == nil {
		return nil, fmt.Errorf("flight service is nil")
	}
	
	// Validate input
	if departureCity == "" || destinationCity == "" {
		return nil, fmt.Errorf("departure city or destination city is empty")
	}
	
	
	cities := fmt.Sprintf("%s_%s", departureCity, destinationCity)
	fmt.Printf("Getting flight data for %s to %s\n", departureCity, destinationCity)

	// Check flight cache
	if fm.FlightService != nil {
		if travelData := fm.FlightService.GetCache(cities); travelData != nil {
			fmt.Printf("Found flight data in cache for %s to %s\n", departureCity, destinationCity)
			return travelData, nil
		}
	}

	// Retrieve airport data from cache or fetch if missing
	var airports *dtos.Trip
	if fm.AirportService != nil {
		airports = fm.AirportService.GetCache(cities)
	}
	
	if airports == nil {
		fmt.Printf("Airport data not found in cache for %s to %s, fetching...\n", departureCity, destinationCity)
		var err error
		
		if fm.AirportService == nil {
			return nil, fmt.Errorf("airport service is nil when trying to get airports")
		}
		
		airports, err = fm.AirportService.GetAirports(departureCity, destinationCity)
		fmt.Printf("GetAirports result: %v, error: %v\n", airports, err)
		
		if err != nil {
			return nil, fmt.Errorf("failed to get airports: %w", err)
		}
		
		// Validate data
		if airports == nil {
			return nil, fmt.Errorf("failed to get airport data for %s and %s", departureCity, destinationCity)
		}
		
		// Check for valid IATA codes
		if airports.DepartureAirport.IATA == "" || airports.DestinationAirport.IATA == "" {
			return nil, fmt.Errorf("missing IATA code for airports: %s, %s", departureCity, destinationCity)
		}
		
		if fm.AirportService != nil {
			fm.AirportService.SetCache(cities, *airports)
			fmt.Printf("Cached airport data for %s to %s\n", departureCity, destinationCity)
		}
	}

	// Fetch flight estimate and update cache
	if fm.FlightService == nil {
		return nil, fmt.Errorf("flight service is nil when trying to get flight estimate")
	}
	
	fmt.Printf("Fetching flight estimate for %s to %s\n", departureCity, destinationCity)
	travelData, err := fm.FlightService.GetFlightEstimate(airports.DepartureAirport, airports.DestinationAirport)
	if err != nil {
		return nil, fmt.Errorf("failed to get flight estimate: %w", err)
	}
	
	// Validate data
	if travelData == nil {
		return nil, fmt.Errorf("received nil response from flight estimate API")
	}
	
	if travelData.Data == nil {
		return nil, fmt.Errorf("received response with nil Data field from flight estimate API")
	}
	
	if travelData.Data.Attributes == nil {
		return nil, fmt.Errorf("received response with nil Data.Attributes field from flight estimate API")
	}
	
	if fm.FlightService != nil {
		fm.FlightService.SetCache(cities, *travelData)
		fmt.Printf("Cached flight data for %s to %s\n", departureCity, destinationCity)
	}
	
	fmt.Printf("Successfully retrieved travel data for %s to %s\n", departureCity, destinationCity)
	return travelData, nil
}

func (fm *FlightManager) AppendTravelDataToMps(mps []*dtos.Mp) ([]*dtos.Mp, error) {
	var errorMessages []string

	if fm == nil {
        return mps, fmt.Errorf("flight manager is nil")
    }
    
    if fm.AirportService == nil {
        return mps, fmt.Errorf("airport service is nil")
    }
    
    if fm.FlightService == nil {
        return mps, fmt.Errorf("flight service is nil")
    }
    
    if mps == nil {
        return nil, fmt.Errorf("mps slice is nil")
    }
    
    for _, mp := range mps {
        if mp == nil {
            fmt.Println("Warning: Encountered nil MP, skipping")
            continue
        }
        
        travelExpenses := mp.Expenses.TravelExpenses
        
		// Check if travelExpenses is nil
        if travelExpenses == nil {
            fmt.Println("Warning: Travel expenses is nil for MP:", mp.MpName)
            continue
        }
        
        for _, expense := range travelExpenses {
            
			// Check if expense is nil
            if expense == nil {
                fmt.Println("Warning: Encountered nil expense, skipping")
                continue
            }
            
            travelLogs := expense.TravelLogs
            
			// Check if travelLogs is nil
            if travelLogs == nil {
                fmt.Println("Warning: Travel logs is nil for expense:", expense.Claim)
				fmt.Println("Warning: Travel logs is nil for mp:", mp.MpName)
                continue
            }

            
            for i := range travelLogs { 

				// Check is KM distance is within threshold
				// Prevent drivable distances from being included as flights
				isDriveDistance, err := fm.AirportService.IsDriveDistance(travelLogs[i].DepartureCity, travelLogs[i].DestinationCity)
				if err != nil {
					fmt.Printf("Error: Failed to check potential drive distance: %s", err)
				}

				if isDriveDistance {
					log := &travelLogs[i]
					log.TransportationMode = constants.GROUND_TRANSPORTATION
				}
                
				// Skip non-air transportation or empty cities
                if travelLogs[i].TransportationMode != constants.AIR_TRANSPORTATION || 
                   travelLogs[i].DepartureCity == constants.EMPTY_STR || 
                   travelLogs[i].DestinationCity == constants.EMPTY_STR ||
				   travelLogs[i].DepartureCity == travelLogs[i].DestinationCity {
                    continue
                }
                
                log := &travelLogs[i]
                
                fmt.Printf("Processing flight: %s to %s\n", log.DepartureCity, log.DestinationCity)
                
                // Get flight data
				// Emissions and distance
                travelData, err := fm.GetFlightData(log.DepartureCity, log.DestinationCity)
                if err != nil {                    
					errorMessages = append(errorMessages, fmt.Sprintf("MP: %s: Error getting flight data for %s to %s: %v\n", mp.MpName, log.DepartureCity, log.DestinationCity, err))
                    continue
                }

                if travelData == nil || travelData.Data == nil || travelData.Data.Attributes == nil {
                    fmt.Printf("Nil data in response for %s to %s\n", log.DepartureCity, log.DestinationCity)
                    continue
                }

                log.TravelData = dtos.TravelData{
                    Distance:      travelData.Data.Attributes.DistanceValue,
                    DistanceUnit:  travelData.Data.Attributes.DistanceUnit,
                    Emissions:     travelData.Data.Attributes.CarbonKilograms,
                    EmissionsUnit: constants.KILOGRAMS,
                }
                
                fmt.Printf("Successfully processed flight data for %s to %s\n", log.DepartureCity, log.DestinationCity)
            }
        }
    }
	extract.WriteFlightErrorsToFile(errorMessages)
    return mps, nil
}

// Write Flight cache to JSON file
// TODO: Refer to this after processing all files
// Will cut costs of CarbonInterface API calls
// Converting this into a cache will save intensive API calls
func (fm *FlightManager) FlightMapToJsonFile(flightCache dtos.FlightCache) error {
	fileName := "flight_cache.json"
	var existingData []dtos.FlightCache

	if data, err := os.ReadFile(fileName); err == nil {
		if len(data) > 0 {
			if err := json.Unmarshal(data, &existingData); err != nil {
				return fmt.Errorf("failed to unmarshal existing data: %s", err)
			}
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to check file: %s", err)
	}

	existingData = append(existingData, flightCache)

	jsonData, err := json.MarshalIndent(existingData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data to json: %s", err)
	}

	if err := os.WriteFile(fileName, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %s", err)
	}

	return nil
}
