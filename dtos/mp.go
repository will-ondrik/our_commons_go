package dtos

// Expendtiture page dtos
type MpWithExpenseCategories struct {
	MpName              Name      `json:"mpName"`
	Constituency        string    `json:"constituency"`
	Caucus              string    `json:"caucus"`
	Salaries            float64   `json:"salaries"`
	TravelExpenses      Category  `json:"travelExpenses"`
	HospitalityExpenses Category  `json:"hospitalityExpenses"`
	ContractExpenses    Category  `json:"contractExpenses"`
	Year                DateRange `json:"year"`
	Quarter             int       `json:"quarter"`
}

type Category struct {
	ExpenseTotal float64 `json:"expenseTotal"`
	Href         string  `json:"href"`
}

type MpExpensesResults struct {
	TravelExpenses      []*TravelExpense
	HospitalityExpenses []*HospitalityExpense
	ContractExpenses    []*ContractExpense
	Errors              []error
}

type Mp struct {
	MpName        Name
	Constituency  string
	Caucus        string
	Year          DateRange
	FiscalQuarter int
	Expenses      Expenses
}

type Expenses struct {
	Totals              ExpenseTotals
	ContractExpenses    []*ContractExpense
	HospitalityExpenses []*HospitalityExpense
	TravelExpenses      []*TravelExpense
}

type ExpenseTotals struct {
	SalariesCost    float64
	ContractCost    float64
	HospitalityCost float64
	TravelCost      float64
	TotalCost       float64
}
