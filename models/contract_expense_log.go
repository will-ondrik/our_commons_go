package models

type ContractExpenseLog struct {
	ID int
	ExpenseId int
	Supplier string
	Description string
	ExpenseDate string
	Amount float64
}