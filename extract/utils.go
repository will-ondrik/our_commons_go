package extract

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHref(cell *goquery.Selection) (string, error) {
	link := cell.Find("a")
	hrefText, _ := link.Attr("href")
	hrefText = strings.TrimSpace(hrefText)
	if hrefText == "" {
		return "", nil
	}

	href := fmt.Sprintf("%s%s", constants.BASE_URL, hrefText)
	return href, nil
}

func FiscalQuarter(text string) (int, error) {
	var quarter int
	var parseErr error
	switch {
	case strings.Contains(text, "First"):
		quarter = 1
	case strings.Contains(text, "Second"):
		quarter = 2
	case strings.Contains(text, "Third"):
		quarter = 3
	case strings.Contains(text, "Fourth"):
		quarter = 4
	default:
		quarter = -1
		parseErr = fmt.Errorf("Failed to extract quarter. Unknown entry")
	}
	return quarter, parseErr
}

func ReportDates(text string) dtos.DateRange {
	splitText := strings.Split(text, "\n")
	tt := strings.Split(splitText[0], " – ")
	dateRangeStr := tt[1]
	dateRangeStr = strings.ReplaceAll(dateRangeStr, "From ", "")
	dates := strings.Split(dateRangeStr, " to ")

	dateRange := dtos.DateRange{
		StartDate: dates[0],
		EndDate:   dates[1],
	}

	return dateRange
}

func ReportYears(dateRange dtos.DateRange) dtos.DateRange {
	startYearStr := dateRange.StartDate
	splitStart := strings.Split(startYearStr, ", ")

	endYearStr := dateRange.EndDate
	splitEnd := strings.Split(endYearStr, ", ")

	return dtos.DateRange{
		StartDate: splitStart[1],
		EndDate:   splitEnd[1],
	}
}

/*
func MpAllExpenses(mp *dtos.MpWithExpenseCategories, b *browser.Browser) dtos.MpExpensesResults {
	var wg sync.WaitGroup

	travelChan := make(chan []*dtos.TravelExpense, 1)
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractTravelExpenses",
			Url:                mp.TravelExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		travelExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			if travelExpenses, ok := travelExpenses.([]*dtos.TravelExpense); ok {
				travelChan <- travelExpenses
			} else {
				errChan <- fmt.Errorf("type assertion failed for hospitality expenses")
			}
		}
	}()

	wg.Add(1)
	hospitalityChan := make(chan []*dtos.HospitalityExpense, 1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractHospitalityExpenses",
			Url:                mp.HospitalityExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		hospitalityExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			if hospitalityExpenses, ok := hospitalityExpenses.([]*dtos.HospitalityExpense); ok {
				hospitalityChan <- hospitalityExpenses
			} else {
				errChan <- fmt.Errorf("type assertion failed for hospitality expenses")
			}
		}
	}()

	wg.Add(1)
	contractChan := make(chan []*dtos.ContractExpense, 1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractContractExpenses",
			Url:                mp.ContractExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		contractExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			contractExpenses, ok := contractExpenses.([]*dtos.ContractExpense)
			if !ok {
				errChan <- err
			} else {
				contractChan <- contractExpenses
			}
		}
	}()

	go func() {
		wg.Wait()
		close(travelChan)
		close(hospitalityChan)
		close(contractChan)
		close(errChan)
	}()

	result := dtos.MpExpensesResults{}

	for travel := range travelChan {
		result.TravelExpenses = travel
	}
	for hospitality := range hospitalityChan {
		result.HospitalityExpenses = hospitality
	}
	for contract := range contractChan {
		result.ContractExpenses = contract
	}
	for err := range errChan {
		result.Errors = append(result.Errors, err)
	}
	return result
}*/
