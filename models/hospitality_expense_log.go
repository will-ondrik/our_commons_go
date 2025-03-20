package models

type HospitalityExpenseLog struct {
	ID int
	EventId int
	ClaimNumber string
	Supplier string
	Amount float64
}