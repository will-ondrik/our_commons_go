package extract

import (
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Mps(doc *goquery.Document) ([]*dtos.MpWithExpenseCategories, error) {
	var mps []*dtos.MpWithExpenseCategories
	var parseErr error

	dataTable := doc.Find("#data-table")
	dataTable.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		mp := &dtos.MpWithExpenseCategories{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {

			text := strings.TrimSpace(cell.Text())
			switch j {
			case 0:
				format.Name(text, mp)
			case 1:
				mp.Constituency = text
			case 2:
				mp.Caucus = text
			case 3:
				salaries, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.Salaries = salaries
			case 4:
				travelExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.TravelExpenses.ExpenseTotal = travelExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
					return
				}
				mp.TravelExpenses.Href = href
			case 5:
				hospitalityExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.HospitalityExpenses.ExpenseTotal = hospitalityExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
					return
				}
				mp.HospitalityExpenses.Href = href
			case 6:
				contractExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.ContractExpenses.ExpenseTotal = contractExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
				}
				mp.ContractExpenses.Href = href
			default:
				fmt.Println("Unexpected column")

			}
		})

		if parseErr == nil {
			mps = append(mps, mp)
		} else {
			log.Printf("\n\nError parsing %s %s: %+v", mp.MpName.FirstName, mp.MpName.LastName, parseErr)
		}
	})

	return mps, parseErr
}
