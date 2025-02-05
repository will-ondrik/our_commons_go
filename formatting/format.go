package format

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Name parses a full name and sets the MP's first and last names.
// It assumes that non-vacant names are in the format "LastName, FirstName"
// and that any honorific appears at the beginning.
func Name(fullName string, mp *dtos.MpWithExpenseCategories) {
	// Handle the special case.
	if fullName == "Vacant" {
		mp.MpName.FirstName = constants.VACANT
		mp.MpName.LastName = constants.VACANT
		return
	}

	var formattedStr string
	// Use HasPrefix (more efficient than Contains+ReplaceAll) assuming honorifics appear at the start.
	if strings.HasPrefix(fullName, "Right Hon. ") {
		formattedStr = fullName[len("Right Hon. "):]
	} else if strings.HasPrefix(fullName, "Hon. ") {
		formattedStr = fullName[len("Hon. "):]
	} else {
		formattedStr = fullName
	}

	// Instead of splitting into a slice, use Index to avoid extra allocations.
	i := strings.Index(formattedStr, ", ")
	if i < 0 {
		// Fallback: if the comma isn't found, assign the whole string as last name.
		mp.MpName.FirstName = ""
		mp.MpName.LastName = formattedStr
	} else {
		mp.MpName.FirstName = formattedStr[i+2:]
		mp.MpName.LastName = formattedStr[:i]
	}
}

// ExpenseToFloat converts a formatted expense string to a float.
// The trimming and replacement operations are chained to reduce intermediate assignments.
func ExpenseToFloat(expenseTotal string) (float64, error) {
	trimmedExpense := strings.ReplaceAll(strings.TrimPrefix(strings.Trim(expenseTotal, "()"), "$"), ",", "")
	expenseFloat, err := strconv.ParseFloat(trimmedExpense, 64)
	if err != nil {
		return -1, fmt.Errorf("Failed to parse string: '%s'\nError: %v", trimmedExpense, err)
	}
	return expenseFloat, nil
}

// StringToDateRange parses a date string into a DateRange.
func StringToDateRange(dateStr string) dtos.DateRange {
	dateStr = strings.TrimSpace(dateStr)
	if dateStr == "" {
		return dtos.DateRange{
			StartDate: "Not Provided",
			EndDate:   "Not Provided",
		}
	}
	dateArr := strings.Split(dateStr, " ")
	return dtos.DateRange{
		StartDate: dateArr[1],
		EndDate:   dateArr[3],
	}
}

// FlightPointsToFloat converts a flight points string to a float.
func FlightPointsToFloat(flightPoints string) (float64, error) {
	pointsFloat, err := strconv.ParseFloat(flightPoints, 64)
	if err != nil {
		return -1, fmt.Errorf("Failed to parse string: %s\nError: %v", flightPoints, err)
	}
	return pointsFloat, nil
}

// TravellerName parses a traveller's name from the format "LastName, FirstName".
func TravellerName(name string) dtos.Name {
	name = strings.TrimSpace(name)
	if strings.Contains(name, "Not Listed") || name == "" {
		return dtos.Name{
			FirstName: "Not Listed",
			LastName:  "Not Listed",
		}
	}
	// Use Index instead of splitting the string.
	i := strings.Index(name, ", ")
	if i < 0 {
		return dtos.Name{
			FirstName: "",
			LastName:  name,
		}
	}
	return dtos.Name{
		FirstName: name[i+2:],
		LastName:  name[:i],
	}
}

// City formats a city name: converts to lower case then capitalizes the first letter.
func City(cityName string) string {
	cityName = strings.ToLower(cityName)
	runes := []rune(cityName)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}
	return string(runes)
}

// Supplier extracts the primary supplier name when both English and French versions are present.
func Supplier(supplier string) string {
	// Instead of using Split, use Index and slicing to avoid extra allocations.
	if i := strings.Index(supplier, " - "); i != -1 {
		supplier = supplier[:i]
	}
	return supplier
}

// EventType cleans up an event type string.
func EventType(text string) string {
	text = strings.TrimSpace(text)
	removeTitle := strings.ReplaceAll(text, "Type of Event", "")
	removeClaim := strings.ReplaceAll(removeTitle, "Claim", "")
	trimmed := strings.TrimSpace(removeClaim)
	return trimmed
}
