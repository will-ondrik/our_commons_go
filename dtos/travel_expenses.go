package dtos

// Travel expense dtos
type TravelExpense struct {
	Claim        string       `json:"claim"`
	Dates        DateRange    `json:"dateRange"`
	TravelLogs   []Traveller  `json:"travelLogs"`
	TravelCosts  TravelCosts  `json:"travelCosts"`
	FlightPoints FlightPoints `json:"flightPoints"`
}
type TravelCosts struct {
	TravelData []TravelData
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
type TravelData struct {
	Distance float64
	DistanceUnit string
	Emissions float64
	EmissionsUnit string
}
