package models

type MpTerm struct {
	ID int
	MpId int
	StartDate string
	EndDate string
	FiscalYear int
	FiscalQuarter int
	Constituency string
	Caucus string
}