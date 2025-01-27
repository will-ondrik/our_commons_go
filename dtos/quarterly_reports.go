package dtos

type AllExpenditureReports struct {
	Reports []ExpenditureReport
}

type ExpenditureReport struct {
	Years     DateRange
	DateRange DateRange
	Quarter   int
	Href      string
}
