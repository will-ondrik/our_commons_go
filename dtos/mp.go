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
	Years                DateRange `json:"year"`
	FiscalYear int `json:"fiscalYear"`
	FiscalQuarter             int       `json:"quarter"`
	Url                 string    `json:"url"`
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
	Years      DateRange
	FiscalYear int
	FiscalQuarter int
	Expenses      Expenses
	Url string
}

type Expenses struct {
	Totals              ExpenseTotals
	ContractExpenses    []*ContractExpense
	ContractExpensesUrl string
	HospitalityExpenses []*HospitalityExpense
	HospitalityExpensesUrl string
	TravelExpenses      []*TravelExpense
	TravelExpensesUrl string
}

type ExpenseTotals struct {
	SalariesCost    float64
	ContractCost    float64
	HospitalityCost float64
	TravelCost      float64
	TotalCost       float64
}
