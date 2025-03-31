package dtos

// TODO: When evaluating Travellers to calculate flight data
// Before calculating TravelData, check that departure and destination city don't match
// There is no details in the data that specify method of transportation
// Add a transportation field to the traveller struct?
type Traveller struct {
	Name            Name   `json:"name"`
	Type            string `json:"type"`
	Purpose         string `json:"purpose"`
	Date            string `json:"date"`
	DepartureCity   string `json:"departureCity"`
	DestinationCity string `json:"destinationCity"`
	TravelData TravelData
	TransportationMode string `json:"transportationMode"`
}