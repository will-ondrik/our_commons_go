package models

type TravelExpenseLog struct {
	ID                        int
	EventId                   int
	TransportationAmount      float64
	AccomodationAmount        float64
	MealsAndIncidentalsAmount float64
	TotalAmount               float64
}