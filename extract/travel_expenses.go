/*package extract

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)
func MpTravelExpenses(doc *goquery.Document) ([]*dtos.TravelExpense, error) {
	fmt.Println("Extracting travel expenses...")
	var mpTravelExpenses []*dtos.TravelExpense
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		if row.HasClass("expenses-main-info") {
			cells := row.Find("td")
			fmt.Printf("Found expenses-main-info row with %d cells\n", cells.Length())

			// Edge case: 6-column row with a hidden row for traveller details
			if cells.Length() == 6 {
				travelExpense := &dtos.TravelExpense{}
				travelExpense.Claim = "Not Provided"

				var departureCity, destinationCity string

				cells.Each(func(j int, cell *goquery.Selection) {
					text := strings.TrimSpace(cell.Text())

					switch j {
					case 0:
						// Populate the full date range (both StartDate and EndDate)
						dates, err := format.StringToDateRange(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.Dates = dates
						fmt.Printf("Parsed date range: %s to %s from text: '%s'\n",
							travelExpense.Dates.StartDate, travelExpense.Dates.EndDate, text)
					case 1:
						departureCity = text
					case 2:
						destinationCity = text
					case 3:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Transportation = cost
					case 4:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Accomodation = cost
					case 5:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.MealsAndIncidentals = cost
					}
				})

				// Process the hidden row for traveller details
				if i+1 < rows.Length() {
					hiddenRow := rows.Eq(i + 1)
					if hiddenRow.Find("table").Length() > 0 {
						hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
							var traveller dtos.Traveller
							text := strings.TrimSpace(nestedRow.Text())

							// For the 6-column case, the nested table has a different structure
							// The first row contains headers, and the second row contains values
							if k == 0 {
								// This is the header row, skip it
								fmt.Printf("Found header row: '%s'\n", text)
							} else if k == 1 {
								// This is the data row
								// Extract purpose and traveller info from the cells
								var purpose, travelerInfo string

								// Find the cells in this row
								nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
									cellText := strings.TrimSpace(cell.Text())
									if l == 0 {
										purpose = cellText
										fmt.Printf("Found purpose: '%s'\n", purpose)
									} else if l == 1 {
										travelerInfo = cellText
										fmt.Printf("Found traveller info: '%s'\n", travelerInfo)
									}
								})

								// Create the traveller
								traveller.Purpose = purpose
								member, name := format.TravellerNameAndType(travelerInfo)
								traveller.Type = member
								traveller.Name = name

								// Use the StartDate from the DateRange as the travel date
								if strings.TrimSpace(travelExpense.Dates.StartDate) != "" {
									traveller.Date = travelExpense.Dates.StartDate
								} else {
									traveller.Date = "1970-01-01" // fallback if for some reason StartDate is empty
								}

								fmt.Printf("Setting date for traveller %s %s: '%s'\n",
									traveller.Name.FirstName, traveller.Name.LastName, traveller.Date)
								traveller.DepartureCity = departureCity
								traveller.DestinationCity = destinationCity

								travelExpense.TravelLogs = append(travelExpense.TravelLogs, traveller)
							}
						})
						i++ // Skip hidden row
					}
				}

				mpTravelExpenses = append(mpTravelExpenses, travelExpense)

			} else {
				// Standard case: a full row with multiple columns
				travelExpense := &dtos.TravelExpense{}
				row.Find("td").Each(func(j int, cell *goquery.Selection) {
					text := strings.TrimSpace(cell.Text())
					switch j {
					case 0:
						travelExpense.Claim = text
					case 1:
						dates, err := format.StringToDateRange(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.Dates = dates
					case 2:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Transportation = cost
					case 3:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Accomodation = cost
					case 4:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.MealsAndIncidentals = cost
					case 5:
						points, err := format.FlightPointsToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.FlightPoints.Regular = points
					case 6:
						points, err := format.FlightPointsToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.FlightPoints.Special = points
					case 7:
						points, err := format.FlightPointsToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.FlightPoints.USA = points
					case 8:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Total = cost
					}
				})

				// Process traveller details from the hidden row
				if i+1 < rows.Length() {
					hiddenRow := rows.Eq(i + 1)
					if hiddenRow.Find("table").Length() > 0 {
						hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
							var traveller dtos.Traveller
							numCols := nestedRow.Find("td").Length()

							nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
								text := strings.TrimSpace(cell.Text())

								if numCols == 4 {
									switch l {
									case 0:
										traveller.Name = format.TravellerName(text)
										fmt.Printf("Processing traveller: %s %s\n",
											traveller.Name.FirstName, traveller.Name.LastName)
									case 1:
										traveller.Type = text
									case 2:
										traveller.Purpose = text
									case 3:
										formattedCity := format.City(text)
										traveller.DepartureCity = formattedCity
										traveller.DestinationCity = formattedCity
									}
									traveller.TransportationMode = constants.GROUND_TRANSPORTATION
									// Use the StartDate from the DateRange as traveller.Date
									traveller.Date = travelExpense.Dates.StartDate
									fmt.Printf("Setting date for 4-column traveller %s %s: '%s'\n",
										traveller.Name.FirstName, traveller.Name.LastName, traveller.Date)
								} else {
									switch l {
									case 0:
										traveller.Name = format.TravellerName(text)
									case 1:
										traveller.Type = text
									case 2:
										traveller.Purpose = text
									case 3:
										traveller.Date = text
									case 4:
										traveller.DepartureCity = text
									case 5:
										traveller.DestinationCity = text
									}
								}
							})

							// Ensure traveller has a valid date before adding to logs
							if strings.TrimSpace(traveller.Date) == "" {
								// If date is empty, use the expense's start date or a fallback
								if strings.TrimSpace(travelExpense.Dates.StartDate) != "" {
									traveller.Date = travelExpense.Dates.StartDate
									fmt.Printf("Using expense start date for traveller: %s\n", traveller.Date)
								} else {
									traveller.Date = "1970-01-01" // fallback date
									fmt.Printf("Using fallback date for traveller: %s\n", traveller.Date)
								}
							}

							fmt.Printf("Adding traveller to logs: %s %s, Date: '%s', Purpose: '%s'\n",
								traveller.Name.FirstName, traveller.Name.LastName,
								traveller.Date, traveller.Purpose)
							travelExpense.TravelLogs = append(travelExpense.TravelLogs, traveller)
						})
						i++ // Skip hidden row
					}
				}

				mpTravelExpenses = append(mpTravelExpenses, travelExpense)
			}
		}
	}

	return mpTravelExpenses, parseErr
}
*/

package extract

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func MpTravelExpenses(doc *goquery.Document) ([]*dtos.TravelExpense, error) {
	fmt.Println("Extracting travel expenses...")
	var mpTravelExpenses []*dtos.TravelExpense
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		if row.HasClass("expenses-main-info") {
			cells := row.Find("td")
			fmt.Printf("Found expenses-main-info row with %d cells\n", cells.Length())
			cells.Each(func(i int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				fmt.Printf("Cell %d: '%s'\n", i, text)
			})
			

			// Skip invalid rows
			if cells.Length() == 0 || strings.Contains(strings.ToLower(cells.Text()), "traveller") || strings.TrimSpace(cells.Eq(0).Text()) == "" {
				fmt.Println("Skipping invalid or header row")
				continue
			}

			// CASE 1: 8-column row
			if cells.Length() == 8 {
				travelExpense := &dtos.TravelExpense{}
				travelExpense.Claim = "Not Provided"
				var departureCity, destinationCity string

				cells.Each(func(j int, cell *goquery.Selection) {
					text := strings.TrimSpace(cell.Text())
					switch j {
					case 0:
						dates, err := format.StringToDateRange(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.Dates = dates
					case 1:
						departureCity = format.CityName(text)
					case 2:
						destinationCity = format.CityName(text)
					case 3:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Transportation = cost
					case 4:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Accomodation = cost
					case 5:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.MealsAndIncidentals = cost
					case 6:
						cost, err := format.ExpenseToFloat(text)
						if err != nil {
							parseErr = err
							return
						}
						travelExpense.TravelCosts.Total = cost
					}
				})

				// Hidden nested table
				if i+1 < rows.Length() {
					hiddenRow := rows.Eq(i + 1)
					if hiddenRow.Find("table").Length() > 0 {
						hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
							if k == 0 || strings.TrimSpace(nestedRow.Text()) == "" {
								return
							}
							var traveller dtos.Traveller
							var purpose, travelerInfo string
							nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
								text := strings.TrimSpace(cell.Text())
								if l == 0 {
									purpose = text
								} else if l == 1 {
									travelerInfo = text
								}
							})
							traveller.Purpose = purpose
							traveller.Type, traveller.Name = format.TravellerNameAndType(travelerInfo)

							// Set date
							if _, err := time.Parse("2006-01-02", travelExpense.Dates.StartDate); err == nil {
								traveller.Date = travelExpense.Dates.StartDate
							} else {
								traveller.Date = "1970-01-01"
							}

							traveller.DepartureCity = format.CityName(departureCity)
							traveller.DestinationCity = format.CityName(destinationCity)

							if IsFlight(traveller.Purpose) {
								traveller.TransportationMode = constants.AIR_TRANSPORTATION
							} else {
								traveller.TransportationMode = constants.GROUND_TRANSPORTATION
							}

							travelExpense.TravelLogs = append(travelExpense.TravelLogs, traveller)
						})
						i++
					}
				}
				fmt.Printf("travel expense: %+v\n\n", travelExpense)
				mpTravelExpenses = append(mpTravelExpenses, travelExpense)

			} else {
				// CASE 2: Full row
				travelExpense := &dtos.TravelExpense{}
				row.Find("td").Each(func(j int, cell *goquery.Selection) {
					text := strings.TrimSpace(cell.Text())
					switch j {
					case 0:
						travelExpense.Claim = text
					case 1:
						dates, err := format.StringToDateRange(text)
						if err == nil {
							travelExpense.Dates = dates
						}
					case 2:
						travelExpense.TravelCosts.Transportation, _ = format.ExpenseToFloat(text)
					case 3:
						travelExpense.TravelCosts.Accomodation, _ = format.ExpenseToFloat(text)
					case 4:
						travelExpense.TravelCosts.MealsAndIncidentals, _ = format.ExpenseToFloat(text)
					case 5:
						travelExpense.FlightPoints.Regular, _ = format.FlightPointsToFloat(text)
					case 6:
						travelExpense.FlightPoints.Special, _ = format.FlightPointsToFloat(text)
					case 7:
						travelExpense.FlightPoints.USA, _ = format.FlightPointsToFloat(text)
					case 8:
						travelExpense.TravelCosts.Total, _ = format.ExpenseToFloat(text)
					}
				})

				// Check if nested row with travellers exists
				if i+1 < rows.Length() {
					hiddenRow := rows.Eq(i + 1)
					if hiddenRow.Find("table").Length() > 0 {
						hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
							if strings.TrimSpace(nestedRow.Text()) == "" {
								return
							}

							var traveller dtos.Traveller
							numCols := nestedRow.Find("td").Length()

							nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
								text := strings.TrimSpace(cell.Text())
								switch numCols {
								case 4:
									switch l {
									case 0:
										traveller.Name = format.TravellerName(text)
									case 1:
										traveller.Type = text
									case 2:
										traveller.Purpose = text
									case 3:
										city := format.City(text)
										traveller.DepartureCity = format.CityName(city)
										traveller.DestinationCity = format.CityName(city)
									}
									if IsFlight(traveller.Purpose) {
										traveller.TransportationMode = constants.AIR_TRANSPORTATION
									} else {
										traveller.TransportationMode = constants.GROUND_TRANSPORTATION
									}
									traveller.Date = travelExpense.Dates.StartDate
								default:
									switch l {
									case 0:
										traveller.Name = format.TravellerName(text)
									case 1:
										traveller.Type = text
									case 2:
										traveller.Purpose = text
									case 3:
										traveller.Date = text
									case 4:
										traveller.DepartureCity = format.CityName(text)
									case 5:
										traveller.DestinationCity = format.CityName(text)
									}
									if IsFlight(traveller.Purpose) {
										traveller.TransportationMode = constants.AIR_TRANSPORTATION
									} else {
										traveller.TransportationMode = constants.GROUND_TRANSPORTATION
									}								
								}
							})

							if traveller.Name.FirstName == "" || traveller.Type == "" {
								return
							}

							if traveller.Date == "" {
								
								isCancelledTrip := format.IsTripCancelled(traveller.Purpose)
								if isCancelledTrip {
									format.HandleCancelledTrip(&traveller, travelExpense.Dates)
								}

								isLoungeVisit := format.IsLoungeVisit(traveller.Purpose)
								if isLoungeVisit {
									format.HandleLoungeVisit(&traveller, travelExpense.Dates)
								}
							}

							travelExpense.TravelLogs = append(travelExpense.TravelLogs, traveller)
						})
						i++
					}
				}
				mpTravelExpenses = append(mpTravelExpenses, travelExpense)
			}
		}
	}

	return mpTravelExpenses, parseErr
}
