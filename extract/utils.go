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
	"time"

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
	fmt.Println("ReportDates input text:", text)
	
	splitText := strings.Split(text, "\n")
	if len(splitText) == 0 {
		return dtos.DateRange{}, fmt.Errorf("invalid text format: no lines found")
	}
	
	fmt.Println("splitText[0]:", splitText[0])
	
	tt := strings.Split(splitText[0], " – ")
	fmt.Println("tt after split:", tt)
	
	if len(tt) < 2 {
		tt = strings.Split(splitText[0], " - ")
		fmt.Println("tt after regular dash split:", tt)
		
		if len(tt) < 2 {
			return dtos.DateRange{}, fmt.Errorf("invalid text format: missing date range separator")
		}
	}
	
	dateRangeStr := tt[1]
	fmt.Println("dateRangeStr:", dateRangeStr)
	
	dateRangeStr = strings.ReplaceAll(dateRangeStr, "From ", "")
	fmt.Println("dateRangeStr after removing 'From ':", dateRangeStr)
	
	dates := strings.Split(dateRangeStr, " to ")
	fmt.Println("dates after split:", dates)
	
	if len(dates) < 2 {
		return dtos.DateRange{}, fmt.Errorf("invalid date range format: missing 'to' separator")
	}

	startDate := strings.TrimSpace(dates[0])
	endDate := strings.TrimSpace(dates[1])
	
	fmt.Println("startDate:", startDate)
	fmt.Println("endDate:", endDate)
	
	// Check if already in Postgres format
	if format.IsPostgresDate(startDate) && format.IsPostgresDate(endDate) {
		fmt.Println("Dates are already in PostgreSQL format")
		return dtos.DateRange{
			StartDate: startDate,
			EndDate: endDate,
		}, nil
	}

	dateRange := dtos.DateRange{
		StartDate: startDate,
		EndDate:   endDate,
	}
	
	fmt.Println("dateRange before conversion:", dateRange)

	formattedDateRange, err := format.ConvertDateFormat(dateRange)
	if err != nil {
		// If conversion fails, return the original dates
		fmt.Println("Conversion failed:", err)
		return dateRange, nil
	}
	
	fmt.Println("formattedDateRange after conversion:", formattedDateRange)

	return formattedDateRange, nil
}

func ReportYear(dateRange dtos.DateRange) (int, error) {
	fmt.Println("ReportYear input dateRange:", dateRange)
	
	// Handle empty date
	if dateRange.StartDate == "" {
		// Default to current year if date is empty
		currentYear := time.Now().Year()
		fmt.Println("Using current year as fallback:", currentYear)
		return currentYear, nil
	}
	
	fmt.Println("StartDate:", dateRange.StartDate)
	
	// Check if the date is in Postgres format
	if format.IsPostgresDate(dateRange.StartDate) {
		startDateSplit := strings.Split(dateRange.StartDate, "-")
		fmt.Println("startDateSplit:", startDateSplit)
		
		if len(startDateSplit) >= 1 {
			// Handle empty year part
			yearStr := strings.TrimSpace(startDateSplit[0])
			fmt.Println("yearStr:", yearStr)
			
			if yearStr != "" {
				year, err := strconv.Atoi(yearStr)
				if err == nil {
					fmt.Println("Extracted year:", year)
					return year, nil
				}
				fmt.Println("Error converting year:", err)
			}
		}
	}
	
	// If year cannot be extracted, try to parse it as a full date
	for _, format := range []string{
		"2006-01-02",
		"January 2, 2006",
		"Jan 2, 2006",
		"2006/01/02",
	} {
		t, err := time.Parse(format, dateRange.StartDate)
		if err == nil {
			year := t.Year()
			fmt.Println("Parsed year using format", format, ":", year)
			return year, nil
		}
	}
	
	// If all fails, use the current year
	currentYear := time.Now().Year()
	fmt.Println("Using current year as last resort:", currentYear)
	return currentYear, nil
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

