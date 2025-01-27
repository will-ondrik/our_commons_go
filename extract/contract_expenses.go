package extract

import (
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func MpContractExpenses(doc *goquery.Document) ([]*dtos.ContractExpense, error) {
	var contractExpenses []*dtos.ContractExpense
	var parseErr error

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		contractExpense := &dtos.ContractExpense{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())

			switch j {
			case 0:
				contractExpense.Supplier = format.Supplier(text)
			case 1:
				contractExpense.Description = text
			case 2:
				contractExpense.Date = text
			case 3:
				expense, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					break
				}
				contractExpense.Total = expense
			}
		})
		contractExpenses = append(contractExpenses, contractExpense)
	})

	return contractExpenses, parseErr
}
