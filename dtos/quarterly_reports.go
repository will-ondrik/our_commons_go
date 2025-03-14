package dtos

type AllExpenditureReports struct {
	Reports []ExpenditureReport
}

type ExpenditureReport struct {
	DateRange DateRange
	FiscalYear int
	FiscalQuarter   int
	Href      string
}
