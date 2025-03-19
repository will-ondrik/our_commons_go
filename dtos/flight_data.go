package dtos

type CarbonInterfaceResponse struct {
	Data *CarbonInterfaceData `json:"data"`
}

type CarbonInterfaceData struct {
	ID         string      `json:"id"`
	Type       string      `json:"type"`
	Attributes *Attributes `json:"attributes"`
}

type Attributes struct {
	Passengers         int       `json:"passengers"`
	Legs              []Flights `json:"legs"`
	DistanceUnit      string    `json:"distance_unit"`
	DistanceValue     float64   `json:"distance_value"`
	EstimatedAt       string    `json:"estimated_at"`
	CarbonGrams       int       `json:"carbon_g"` 
	CarbonPounds      float64   `json:"carbon_lb"`
	CarbonKilograms   float64   `json:"carbon_kg"`
	CarbonMetricTonnes float64  `json:"carbon_mt"`
}

type Flights struct {
	DepartureAirport   string `json:"departure_airport"`
	DestinationAirport string `json:"destination_airport"`
}

type CarbonInterfaceRequest struct {
	Type       string    `json:"type"`
	Passengers int       `json:"passengers"`
	Legs       []Flights `json:"legs"`
}

type FlightCache map[string]CarbonInterfaceResponse
