package dtos

// Travel expense dtos
type TravelExpenses struct {
	Claim        string       `json:"claim"`
	Dates        DateRange    `json:"dateRange"`
	TravelLogs   []Traveller  `json:"travelLogs"`
	TravelCosts  TravelCosts  `json:"travelCosts"`
	FlightPoints FlightPoints `json:"flightPoints"`
}
type TravelCosts struct {
	Transportation      float64 `json:"transportation"`
	Accomodation        float64 `json:"accomodation"`
	MealsAndIncidentals float64 `json:"mealsAndIncidentals"`
	Total               float64 `json:"total"`
}
type FlightPoints struct {
	Regular float64 `json:"regular"`
	Special float64 `json:"special"`
	USA     float64 `json:"usa"`
}
