package models

import "etl_our_commons/dtos"

type Mp struct {
	ID int
	MpName        dtos.Name
	Constituency  string
	Caucus        string
	Years      dtos.DateRange
	FiscalYear int
	FiscalQuarter int
	Expenses      dtos.Expenses
}

