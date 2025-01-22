package extract

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"log"

	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func Mps(doc *goquery.Document) ([]*dtos.MpWithExpenseCategories, error) {
	var mps []*dtos.MpWithExpenseCategories
	var parseErr error

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {

		mp := &dtos.MpWithExpenseCategories{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {

			text := strings.TrimSpace(cell.Text())
			switch j {
			case 0:
				format.Name(text, mp)
			case 1:
				mp.Constituency = text
			case 2:
				mp.Caucus = text
			case 3:
				salaries, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.Salaries = salaries
			case 4:
				travelExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.TravelExpenses.ExpenseTotal = travelExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
					return
				}
				mp.TravelExpenses.Href = href
			case 5:
				hospitalityExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.HospitalityExpenses.ExpenseTotal = hospitalityExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
					return
				}
				mp.HospitalityExpenses.Href = href
			case 6:
				contractExpenses, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				mp.ContractExpenses.ExpenseTotal = contractExpenses

				href, err := getHref(cell)
				if err != nil {
					parseErr = err
				}
				mp.ContractExpenses.Href = href
			default:
				fmt.Println("Unexpected column")

			}
		})

		if parseErr == nil {
			mps = append(mps, mp)
		} else {
			log.Printf("\n\nError parsing %s %s: %+v", mp.MpName.FirstName, mp.MpName.LastName, parseErr)
		}
	})

	return mps, parseErr
}

func getHref(cell *goquery.Selection) (string, error) {
	link := cell.Find("a")
	hrefText, _ := link.Attr("href")

	href := fmt.Sprintf("%s%s", constants.BASE_URL, hrefText)
	return href, nil
}

func MpTravelExpenses(doc *goquery.Document) ([]*dtos.TravelExpenses, error) {
	fmt.Println("Extracting travel expenses...")
	var mpTravelExpenses []*dtos.TravelExpenses
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		// Process the visible row
		if row.HasClass("expenses-main-info") {
			var travelExpenses dtos.TravelExpenses

			row.Find("td").Each(func(j int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				switch j {
				case 0:
					travelExpenses.Claim = text
				case 1:
					travelExpenses.Dates = format.StringToDateRange(text)
				case 2:
					transportationCosts, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.TravelCosts.Transportation = transportationCosts
				case 3:
					accommodationCosts, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.TravelCosts.Accomodation = accommodationCosts
				case 4:
					mealsAndIncidentals, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.TravelCosts.MealsAndIncidentals = mealsAndIncidentals
				case 5:
					points, err := format.FlightPointsToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.FlightPoints.Regular = points
				case 6:
					points, err := format.FlightPointsToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.FlightPoints.Special = points
				case 7:
					points, err := format.FlightPointsToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.FlightPoints.USA = points
				case 8:
					total, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					travelExpenses.TravelCosts.Total = total
				}
			})

			// Check for the hidden row
			if i+1 < rows.Length() {
				hiddenRow := rows.Eq(i + 1)
				if hiddenRow.Find("table").Length() > 0 {
					hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
						var traveller dtos.Traveller
						numCols := nestedRow.Find("td").Length()

						nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
							text := strings.TrimSpace(cell.Text())
							// Case for transportation without flights
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
							} else { // Normal case
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

						travelExpenses.TravelLogs = append(travelExpenses.TravelLogs, traveller)
					})
					i++
				}
			}

			mpTravelExpenses = append(mpTravelExpenses, &travelExpenses)
		}
	}

	return mpTravelExpenses, parseErr
}

func MpHospitalityExpenses() {

}

func MpContractExpenses() {

}

func MpAllExpenses() {

}
