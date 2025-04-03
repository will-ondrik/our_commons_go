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
		
		// Check if the name contains a comma before splitting
		if strings.Contains(formattedStr, ", ") {
			names := strings.Split(formattedStr, ", ")
			firstName = strings.TrimSpace(names[1])
			lastName = strings.TrimSpace(names[0])
		} else {
			// Handle case where name doesn't have a comma
			// Set the entire name as LastName
			firstName = ""
			lastName = formattedStr
		}
	}
	mp.MpName.FirstName = firstName
	mp.MpName.LastName = lastName
}

// TODO: Refactor - some of this is unnecessary
func TravellerNameAndType(travellerDetails string) (string, dtos.Name) {
	fmt.Println("Traveller details string:", travellerDetails)
	// Check if it contains a name in parentheses
	if strings.Contains(travellerDetails, "(") && strings.Contains(travellerDetails, ")") {
		split := strings.SplitN(travellerDetails, "(", 2)
		travellerType := strings.TrimSpace(split[0])

		// Clean and extract name
		nameRaw := strings.TrimSpace(strings.TrimRight(split[1], ")"))
		nameSlice := strings.Split(nameRaw, ", ")

		if len(nameSlice) == 2 {
			return travellerType, dtos.Name{
				FirstName: strings.TrimSpace(nameSlice[1]),
				LastName:  strings.TrimSpace(nameSlice[0]),
			}
		}

		// Fallback if name doesn't split as expected
		return travellerType, dtos.Name{
			FirstName: "",
			LastName:  nameRaw,
		}
	}

	// No parentheses - treat whole string as name or type
	return "Unknown", dtos.Name{
		FirstName: "",
		LastName:  strings.TrimSpace(travellerDetails),
	}
}



func ExpenseToFloat(expenseTotal string) (float64, error) {
	// Handle empty strings
	expenseTotal = strings.TrimSpace(expenseTotal)
	if expenseTotal == "" {
		return 0, nil
	}
	
	trimmedExpense := strings.Trim(expenseTotal, "()")
	trimmedExpense = strings.TrimPrefix(trimmedExpense, "$")
	trimmedExpense = strings.ReplaceAll(trimmedExpense, ",", "")

	expenseFloat, err := strconv.ParseFloat(trimmedExpense, 64)
	if err != nil {
		return 0, fmt.Errorf("Failed to parse string: '%s'\nError: %v", trimmedExpense, err)
	}

	return expenseFloat, nil
}

func StringToDateRange(dateStr string) (dtos.DateRange, error) {
	fmt.Println("date string:", dateStr)
	dateStr = strings.TrimSpace(dateStr)

	if dateStr == "" {
		return dtos.DateRange{
			StartDate: "Not Provided",
			EndDate:   "Not Provided",
		}, nil
	}

	// Normalize spacing and remove optional "From"
	dateStr = strings.ReplaceAll(dateStr, "From ", "")
	dateStr = strings.ReplaceAll(dateStr, "from ", "")
	dateStr = strings.Join(strings.Fields(dateStr), " ") // Collapse spaces

	// Safe "to" range check
	if strings.Contains(dateStr, " to ") {
		dateArr := strings.Split(dateStr, " to ")
		if len(dateArr) == 2 {
			startDate := strings.TrimSpace(dateArr[0])
			endDate := strings.TrimSpace(dateArr[1])

			// Check if already in PSQL format
			if IsPostgresDate(startDate) && IsPostgresDate(endDate) {
				return dtos.DateRange{StartDate: startDate, EndDate: endDate}, nil
			}

			// Try conversion
			dateRange := dtos.DateRange{StartDate: startDate, EndDate: endDate}
			formatted, err := ConvertDateFormat(dateRange)
			if err != nil {
				return dateRange, nil
			}
			return formatted, nil
		}
	}

	// Handle single PSQL-style date
	if IsPostgresDate(dateStr) {
		return dtos.DateRange{
			StartDate: dateStr,
			EndDate:   dateStr,
		}, nil
	}

	// Handle malformed but dash-containing ranges
	if strings.Contains(dateStr, "-") {
		parts := strings.Fields(dateStr)
		if len(parts) == 2 {
			return dtos.DateRange{
				StartDate: parts[0],
				EndDate:   parts[1],
			}, nil
		}
	}

	// Final fallback
	return dtos.DateRange{
		StartDate: dateStr,
		EndDate:   dateStr,
	}, nil
}



func FlightPointsToFloat(flightPoints string) (float64, error) {
	
	// Handle empty strings
	flightPoints = strings.TrimSpace(flightPoints)
	if flightPoints == "" {
		return 0, nil
	}
	
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
	
	// Check if the name contains a comma before splitting
	if strings.Contains(name, ", ") {
		nameArr := strings.Split(name, ", ")
		return dtos.Name{
			FirstName: nameArr[1],
			LastName:  nameArr[0],
		}
	} else {
		// Handle case where name doesn't have a comma
		// Set the entire name as LastName
		return dtos.Name{
			FirstName: "",
			LastName:  name,
		}
	}
}

// Special case in Travel expenses
// The city listed is always in all caps
// Format for proper punctuation
func City(cityName string) string {
	cityName = strings.TrimSpace(cityName)
	if cityName == "" {
		return "Not Provided"
	}
	
	cityName = strings.ToLower(cityName)
	runes := []rune(cityName)
	if len(runes) > 0 {
		runes[0] = unicode.ToUpper(runes[0])
	}

	return string(runes)
}

func Supplier(supplier string) string {
	if strings.ContainsAny(supplier, " - ") {
		// English and French version present
		versions := strings.Split(supplier, " - ")
		if len(versions) > 0 {
			supplier = versions[0]
		}
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

func IsTripCancelled(travelPurpose string) bool {
	if strings.Contains(travelPurpose, "cancel") {
		return true
	}
	return false
}

func HandleCancelledTrip(traveller *dtos.Traveller, dates dtos.DateRange) {
	// Set date to empty string if dates.StartDate is "Not Provided"
	if dates.StartDate == "Not Provided" {
		traveller.Date = ""  // This will be converted to NULL in the database
	} else {
		traveller.Date = dates.StartDate
	}
	traveller.DepartureCity = "Not Listed"
	traveller.DestinationCity = "Not Listed"
}

func IsLoungeVisit(travelPurpose string) bool {
	if strings.Contains(travelPurpose, "Maple Leaf Lounge") {
		return true
	}
	return false
}

func HandleLoungeVisit(traveller *dtos.Traveller, dates dtos.DateRange) {
	// Set date to empty string if dates.StartDate is "Not Provided"
	if dates.StartDate == "Not Provided" {
		traveller.Date = ""  // This will be converted to NULL in the database
	} else {
		traveller.Date = dates.StartDate
	}
	traveller.DepartureCity = "Not Listed"
	traveller.DestinationCity = "Not Listed"
}

// Remove extra parenthese and their corresponding text
// This will prevent airport search failure
func CityName(city string) string {

	if strings.Contains(city, "montréal, qc") || strings.Contains(city, "montréal") {
		return "Montreal"
	}
	
	var formatted []rune
	insideParen := false

	for _, v := range city {
		if v == '(' {
			insideParen = true
			continue
		}

		if v == ')' {
			insideParen = false
			continue
		}

		if !insideParen {
			formatted = append(formatted, v)
		}
	}

	return string(formatted)
}

func IsPostgresDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
