package dtos

// General Expense dtos
type DateRange struct {
	StartDate string `json:"startDate"`
	End       string `json:"endDate"`
}

type Name struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
}
