package dtos

type MpDetailsUrl struct {
	FirstName string
	LastName string
	Url string
}

// details url = https:///www.ourcommons.ca/Members/en/ziad-aboultaif(89165) --> this is appended to the base url: /Members/en/ziad-aboultaif(89165)

type MpDetails struct {
	FirstName string
	LastName string
	Province string
	Contact ContactInfo
	Roles []MpRoles
}

type MpRoles struct {
	Current []Role
	ExecuiveCommittees []Role
}

type Role struct {
	Abbrev string
	FullName string
}

type ContactInfo struct {
	Email string
	Website string
	HillOffice Address
	ConstituencyOffice Address
}

type Address struct {
	Name string
	StreetName string
	City string
	Province string
	PostalCode string
	Phone string
	Fax string
	Info string
}