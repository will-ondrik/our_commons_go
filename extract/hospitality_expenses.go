package extract

import (
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func MpHospitalityExpenses(doc *goquery.Document) ([]*dtos.HospitalityExpense, error) {
	rows := doc.Find("tbody tr")
	hospitalityExpenses := make([]*dtos.HospitalityExpense, 0, rows.Length())
	var parseErr error

	nRows := rows.Length()
	for i := 0; i < nRows; i++ {
		row := rows.Eq(i)

		if !row.HasClass("expenses-main-info") {
			continue
		}

		hospitalityExpense := &dtos.HospitalityExpense{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			switch j {
			case 0:
				hospitalityExpense.Date = text
			case 1:
				hospitalityExpense.Location = text
			case 2:
				attendance, err := strconv.Atoi(text)
				if err != nil {
					parseErr = err
					return
				}
				hospitalityExpense.Attendance = attendance
			case 3:
				hospitalityExpense.Purpose = text
			case 4:
				cost, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				hospitalityExpense.TotalCost = cost
			}
		})

		// Process hidden data in the next row
		if i+1 < nRows {
			hiddenRow := rows.Eq(i + 1)
			// Retrieve Event type
			eventType := hiddenRow.Find(".col-md-4").Text()
			hospitalityExpense.Event.Type = format.EventType(eventType)

			var expenseLogs []dtos.ExpenseLogs
			hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
				var expenseEntry dtos.ExpenseLogs
				nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
					text := strings.TrimSpace(cell.Text())
					switch l {
					case 0:
						expenseEntry.Claim = text
					case 1:
						expenseEntry.Supplier = text
					case 2:
						expense, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						expenseEntry.Cost = expense
					}
				})
				// Only add valid entries.
				if expenseEntry.Claim != "" && expenseEntry.Cost != 0 && expenseEntry.Supplier != "" {
					expenseLogs = append(expenseLogs, expenseEntry)
				}
			})
			hospitalityExpense.Event.ExpenseLogs = expenseLogs
		}
		hospitalityExpenses = append(hospitalityExpenses, hospitalityExpense)
	}
	return hospitalityExpenses, parseErr
}
