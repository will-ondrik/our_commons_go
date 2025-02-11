package dtos

import "github.com/PuerkitoBio/goquery"

// Store MP-related html
type MpHtml struct {
	Info        MpWithExpenseCategories
	Contract    html
	Hospitality html
	Travel      html
}

type html *goquery.Document
