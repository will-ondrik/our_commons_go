package extract

import (
	//"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"strconv"

	"strings"

	//"sync"

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

func ReportDates(text string) (dtos.DateRange, error) {
	splitText := strings.Split(text, "\n")
	tt := strings.Split(splitText[0], " – ")
	dateRangeStr := tt[1]
	dateRangeStr = strings.ReplaceAll(dateRangeStr, "From ", "")
	dates := strings.Split(dateRangeStr, " to ")

	dateRange := dtos.DateRange{
		StartDate: dates[0],
		EndDate:   dates[1],
	}

	formattedDateRange, err := format.ConvertDateFormat(dateRange)
	if err != nil {
		return dtos.DateRange{}, err
	}

	return formattedDateRange, nil
}

func ReportYear(dateRange dtos.DateRange) (int, error) {
	startDateSplit := strings.Split(dateRange.StartDate, "-")

	year, err := strconv.Atoi(startDateSplit[0])
	if err != nil {
		return 0, fmt.Errorf("failed to convert year")
	}

	
	return year, nil
}