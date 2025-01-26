package dtos

type AllExpenditureReports struct {
	Reports []ExpenditureReport
}

type ExpenditureReport struct {
	Year    string
	Quarter float64
	Href    string
}
