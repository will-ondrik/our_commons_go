package dtos

type MpWithExpenseCategories struct {
	MpName              MpName
	Constituency        string
	Caucus              string
	Salaries            float64
	TravelExpenses      Category
	HospitalityExpenses Category
	ContractExpenses    Category
}

type Category struct {
	Name         string
	ExpenseTotal float64
	Href         string
}

type MpName struct {
	FirstName string
	LastName  string
}
