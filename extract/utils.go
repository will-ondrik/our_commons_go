package extract

import (
	//"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"os"
	"slices"
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
	if len(splitText) == 0 {
		return dtos.DateRange{}, fmt.Errorf("invalid text format: no lines found")
	}
	
	tt := strings.Split(splitText[0], " – ")
	if len(tt) < 2 {
		return dtos.DateRange{}, fmt.Errorf("invalid text format: missing date range separator")
	}
	
	dateRangeStr := tt[1]
	dateRangeStr = strings.ReplaceAll(dateRangeStr, "From ", "")
	dates := strings.Split(dateRangeStr, " to ")
	
	if len(dates) < 2 {
		return dtos.DateRange{}, fmt.Errorf("invalid date range format: missing 'to' separator")
	}

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
	// Handle empty date
	if dateRange.StartDate == "" {
		return 0, fmt.Errorf("empty date string")
	}
	
	startDateSplit := strings.Split(dateRange.StartDate, "-")
	
	if len(startDateSplit) == 0 {
		return 0, fmt.Errorf("invalid date format: missing year")
	}
	
	// Handle empty year part
	yearStr := strings.TrimSpace(startDateSplit[0])
	if yearStr == "" {
		return 0, fmt.Errorf("empty year part in date")
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, fmt.Errorf("failed to convert year: %v", err)
	}
	
	return year, nil
}



func IsFlight(travelPurpose string) bool {
	// Handle empty strings
	if travelPurpose == "" {
		return false
	}

	travelPurpose = strings.ToLower(travelPurpose)

	// to attend a national caucus meeting
	// to attend a regional or provincial caucus meeting
	// Attending event with Member (type: Employee)
	// // Need to compare cost to ensure its a flight
		// There may be multiple cities in close proximity
	// unite the family with the Member
	// travel to/from constituency and Ottawa

	// Potentials:
		// to attend training
		// to attend meetings with stakeholders about business of the House
		// to support a parliamentary association
		// to attend language training
		//

	return slices.Contains(constants.FLIGHT_KEYWORDS, travelPurpose)
}

func WriteFlightErrorsToFile(errorMsgs []string) error {
	if len(errorMsgs) == 0 {
		fmt.Println("No errors to write.")
		return nil
	}

	content := strings.Join(errorMsgs, "\n")
	filePath := "sample.txt"

	// Confirm where we are writing from
	wd, _ := os.Getwd()
	fmt.Printf("Current working directory: %s\n", wd)
	fmt.Printf("Writing to: %s\n", filePath)

	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		// Crash hard with useful context
		return fmt.Errorf("failed to write to %s: %w", filePath, err)
	}

	fmt.Printf("Successfully wrote %d error(s) to %s\n", len(errorMsgs), filePath)
	return nil
}
