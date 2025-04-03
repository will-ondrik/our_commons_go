package extract

import (
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func MpHospitalityExpenses(doc *goquery.Document) ([]*dtos.HospitalityExpense, error) {
	var hospitalityExpenses []*dtos.HospitalityExpense
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		if !row.HasClass("expenses-main-info") {
			continue
		}

		fmt.Println("Hospitality expenses outer row: ", row.Length())
		fmt.Println("outer row text:", row.Text())
		

		hospitalityExpense := &dtos.HospitalityExpense{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			switch j {
			case 0:
				hospitalityExpense.Date = text
			case 1:
				hospitalityExpense.Location = text
			case 2:
				text = strings.TrimSpace(text)
				if text == "" {
					hospitalityExpense.Attendance = 0 // Default value for empty attendance
				} else {
					attendance, err := strconv.Atoi(text)
					if err != nil {
						parseErr = err
						return
					}
					hospitalityExpense.Attendance = attendance
				}
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

		// Process hidden data
		hiddenRow := rows.Eq(i + 1)
		fmt.Println("Hospitality inner row length: ", hiddenRow.Length())
		fmt.Println("inner row text:", hiddenRow.Text())
		

		// Retrieve Event type
		eventType := hiddenRow.Find(".col-md-4").Text()
		hospitalityExpense.Event.Type = format.EventType(eventType)

		var expenseLogs []dtos.ExpenseLogs
		hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
			fmt.Println("Hospitality inner row length: ", nestedRow.Length())
			fmt.Println("inner row text:", nestedRow.Text())
		
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
			if expenseEntry.Claim == "" || expenseEntry.Cost == 0 || expenseEntry.Supplier == "" {
				return
			}
			expenseLogs = append(expenseLogs, expenseEntry)
		})
		hospitalityExpense.Event.ExpenseLogs = expenseLogs
		hospitalityExpenses = append(hospitalityExpenses, hospitalityExpense)
	}
	WriteErrs(hospitalityExpenses)
	return hospitalityExpenses, parseErr
}

func WriteErrs(expenses []*dtos.HospitalityExpense) error {
	if len(expenses) == 0 {
		fmt.Println("No expenses to write.")
		return nil
	}

	var expenseStrings []string
	for _, expense := range expenses {
		expenseStr := fmt.Sprintf("Date: %s, Location: %s, Purpose: %s, Total Cost: %.2f", 
			expense.Date, expense.Location, expense.Purpose, expense.TotalCost)
		expenseStrings = append(expenseStrings, expenseStr)
	}

	content := strings.Join(expenseStrings, "\n")
	filePath := "hospitality.txt"

	wd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", wd)
	fmt.Printf("Writing to: %s\n", filePath)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %w", filePath, err)
	}

	fmt.Printf("Successfully wrote %d expense(s) to %s\n", len(expenses), filePath)
	return nil
}
