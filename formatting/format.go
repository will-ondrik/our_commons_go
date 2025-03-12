package format

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"
)

func Name(fullName string, mp *dtos.MpWithExpenseCategories) {
	var firstName string
	var lastName string

	switch fullName {

	case "Vacant":
		firstName = constants.VACANT
		lastName = constants.VACANT

	default:
		var formattedStr string

		if strings.Contains(fullName, "Right Hon. ") {
			formattedStr = strings.ReplaceAll(fullName, "Right Hon. ", "")
		} else if strings.Contains(fullName, "Hon. ") {
			formattedStr = strings.ReplaceAll(fullName, "Hon. ", "")
		} else {
			formattedStr = fullName
		}
		names := strings.Split(formattedStr, ", ")

		firstName = names[1]
		lastName = names[0]
	}
	mp.MpName.FirstName = firstName
	mp.MpName.LastName = lastName
}

func ExpenseToFloat(expenseTotal string) (float64, error) {
	trimmedExpense := strings.Trim(expenseTotal, "()")
	trimmedExpense = strings.TrimPrefix(trimmedExpense, "$")
	trimmedExpense = strings.ReplaceAll(trimmedExpense, ",", "")

	expenseFloat, err := strconv.ParseFloat(trimmedExpense, 64)
	if err != nil {
		return -1, fmt.Errorf("Failed to parse string: '%s'\nError: %v", trimmedExpense, err)
	}

	return expenseFloat, nil
}

func StringToDateRange(dateStr string) (dtos.DateRange, error) {
	dateStr = strings.TrimSpace(dateStr)

	if dateStr == "" {
		return dtos.DateRange{
			StartDate: "Not Provided",
			EndDate:   "Not Provided",
		}, nil
	}

	dateArr := strings.Split(dateStr, " ")
	dateRange := dtos.DateRange{
		StartDate: dateArr[1],
		EndDate: dateArr[3],
	}

	formattedDateRange, err := ConvertDateFormat(dateRange)
	if err != nil {
		return dateRange, nil
	}

	fmt.Println("Formatted date:", formattedDateRange)

	return formattedDateRange, nil
}

func FlightPointsToFloat(flightPoints string) (float64, error) {
	pointsFloat, err := strconv.ParseFloat(flightPoints, 64)
	if err != nil {
		return -1, fmt.Errorf("Failed to parse string: %s\nError: %v", flightPoints, err)
	}

	return pointsFloat, nil
}

func TravellerName(name string) dtos.Name {
	name = strings.TrimSpace(name)
	name = strings.ReplaceAll(name, "Right Hon. ", "")
	name = strings.ReplaceAll(name, "Hon. ", "")

	if strings.Contains(name, "Not Listed") || name == "" {
		return dtos.Name{
			FirstName: "Not Listed",
			LastName:  "Not Listed",
		}
	}
	nameArr := strings.Split(name, ", ")

	return dtos.Name{
		FirstName: nameArr[1],
		LastName:  nameArr[0],
	}
}

// Special case in Travel expenses
// The city listed is always in all caps
// Format for proper punctuation
func City(cityName string) string {
	cityName = strings.ToLower(cityName)
	runes := []rune(cityName)
	runes[0] = unicode.ToUpper(runes[0])

	return string(runes)
}

func Supplier(supplier string) string {
	if strings.ContainsAny(supplier, " - ") {
		// English and French version present
		versions := strings.Split(supplier, " - ")
		supplier = versions[0]
	}
	return supplier
}

func EventType(text string) string {
	text = strings.TrimSpace(text)
	removeTitle := strings.ReplaceAll(text, "Type of Event", "")
	removeClaim := strings.ReplaceAll(removeTitle, "Claim", "")
	trimmed := strings.TrimSpace(removeClaim)

	return trimmed
}

// Convert date string into PSQL format
func ConvertDateFormat(dateRange dtos.DateRange) (dtos.DateRange, error) {

	fmt.Println("Convert date input:", dateRange)
	// Must use same date structure in the new format
	// This ensures that the input format can be formatted correctly
	const inputFormat = "January 2, 2006"
	const newFormat = "2006-01-02"

	formattedStartDate, err := time.Parse(inputFormat, dateRange.StartDate)
	if err != nil {
		return dtos.DateRange{}, err
	}

	formattedEndDate, err := time.Parse(inputFormat, dateRange.EndDate)
	if err != nil {
		return dtos.DateRange{}, err
	}

	fmt.Println("New start:", formattedStartDate, "New end:", formattedEndDate)
	return dtos.DateRange{
		StartDate: formattedStartDate.Format(newFormat),
		EndDate: formattedEndDate.Format(newFormat),
	}, nil
}