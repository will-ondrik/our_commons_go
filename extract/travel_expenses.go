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
