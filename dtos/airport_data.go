package dtos

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

type AirportCache map[string]Trip