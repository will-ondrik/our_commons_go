package dtos

type HospitalityExpense struct {
	Date       string
	Location   string
	Attendance int
	Purpose    string
	Event      Event
	TotalCost  float64
}
type Event struct {
	Type        string
	ExpenseLogs []ExpenseLogs
}
type ExpenseLogs struct {
	Claim    string
	Supplier string
	Cost     float64
}
