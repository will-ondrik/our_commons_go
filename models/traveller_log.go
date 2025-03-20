package models

type TravellerLog struct {
	ID int
	EventId int
	TravellerId int
	TravelDate string
	Purpose string
	DepartureCity string
	DestinationCity string
	TranportationMode string
}