package models

type HospitalityEvent struct {
	ID int
	ExpenseId int
	ExpenseDate string
	Location string
	Purpose string
	Amount float64
}