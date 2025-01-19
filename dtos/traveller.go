package dtos

type Traveller struct {
	Name            Name   `json:"name"`
	Type            string `json:"type"`
	Purpose         string `json:"purpose"`
	Date            string `json:"date"`
	DepartureCity   string `json:"departureCity"`
	DestinationCity string `json:"destinationCity"`
}
