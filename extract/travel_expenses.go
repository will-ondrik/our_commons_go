package extract

import (
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func MpTravelExpenses(doc *goquery.Document) ([]*dtos.TravelExpense, error) {
	fmt.Println("Extracting travel expenses...")
	rows := doc.Find("tbody tr")
	mpTravelExpenses := make([]*dtos.TravelExpense, 0, rows.Length())
	var parseErr error
	nRows := rows.Length()

	for i := 0; i < nRows; i++ {
		row := rows.Eq(i)

		// Process only rows with the main expense info.
		if row.HasClass("expenses-main-info") {
			travelExpense := &dtos.TravelExpense{}
			row.Find("td").Each(func(j int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				switch j {
				case 0:
					travelExpense.Claim = text
				case 1:
					travelExpense.Dates = format.StringToDateRange(text)
				case 2:
					transportationCosts, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpense.TravelCosts.Transportation = transportationCosts
				case 3:
					accommodationCosts, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpense.TravelCosts.Accomodation = accommodationCosts
				case 4:
					mealsAndIncidentals, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpense.TravelCosts.MealsAndIncidentals = mealsAndIncidentals
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
					total, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpense.TravelCosts.Total = total
				}
			})

			// Process hidden travel log data if available.
			if i+1 < nRows {
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
								case 1:
									traveller.Type = text
								case 2:
									traveller.Purpose = text
								case 3:
									formattedCity := format.City(text)
									traveller.DepartureCity = formattedCity
									traveller.DestinationCity = formattedCity
								}
								traveller.Date = "Not Provided"
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
						travelExpense.TravelLogs = append(travelExpense.TravelLogs, traveller)
					})
					// Skip the hidden row since it has been processed.
					i++
				}
			}

			mpTravelExpenses = append(mpTravelExpenses, travelExpense)
		}
	}

	return mpTravelExpenses, parseErr
}
