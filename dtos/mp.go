package dtos

// Expendtiture page dtos
type MpWithExpenseCategories struct {
	MpName              Name     `json:"mpName"`
	Constituency        string   `json:"constituency"`
	Caucus              string   `json:"caucus"`
	Salaries            float64  `json:"salaries"`
	TravelExpenses      Category `json:"travelExpenses"`
	HospitalityExpenses Category `json:"hospitalityExpenses"`
	ContractExpenses    Category `json:"contractExpenses"`
	Year                string   `json:"year"`
	Quarter             int      `json:"quarter"`
}

type Category struct {
	Name         string  `json:"name"`
	ExpenseTotal float64 `json:"expenseTotal"`
	Href         string  `json:"href"`
}
