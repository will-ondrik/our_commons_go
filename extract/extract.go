package extract

import (
	//"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	format "etl_our_commons/formatting"
	"fmt"
	"log"

	"strconv"
	"strings"

	//"sync"

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
	hrefText = strings.TrimSpace(hrefText)
	if hrefText == "" {
		return "", nil
	}

	href := fmt.Sprintf("%s%s", constants.BASE_URL, hrefText)
	return href, nil
}

func MpTravelExpenses(doc *goquery.Document) ([]*dtos.TravelExpense, error) {
	fmt.Println("Extracting travel expenses...")
	var mpTravelExpenses []*dtos.TravelExpense
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		// Process the visible row
		if row.HasClass("expenses-main-info") {
			travelExpense := &dtos.TravelExpense{}
			row.Find("td").Each(func(j int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				fmt.Println("Index: ", j)
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

								// Normal case
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
					i++
				}
			}

			mpTravelExpenses = append(mpTravelExpenses, travelExpense)
		}
	}

	return mpTravelExpenses, parseErr
}

func MpHospitalityExpenses(doc *goquery.Document) ([]*dtos.HospitalityExpense, error) {
	var hospitalityExpenses []*dtos.HospitalityExpense
	var parseErr error

	rows := doc.Find("tbody tr")
	for i := 0; i < rows.Length(); i++ {
		row := rows.Eq(i)

		if !row.HasClass("expenses-main-info") {
			continue
		}

		hospitalityExpense := &dtos.HospitalityExpense{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())
			switch j {
			case 0:
				hospitalityExpense.Date = text
			case 1:
				hospitalityExpense.Location = text
			case 2:
				attendance, err := strconv.Atoi(text)
				if err != nil {
					parseErr = err
					return
				}
				hospitalityExpense.Attendance = attendance
			case 3:
				hospitalityExpense.Purpose = text
			case 4:
				cost, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					return
				}
				hospitalityExpense.TotalCost = cost
			}
		})

		// Process hidden data
		hiddenRow := rows.Eq(i + 1)

		// Retrieve Event type
		eventType := hiddenRow.Find(".col-md-4").Text()
		hospitalityExpense.Event.Type = format.EventType(eventType)

		var expenseLogs []dtos.ExpenseLogs
		hiddenRow.Find("table tbody tr").Each(func(k int, nestedRow *goquery.Selection) {
			var expenseEntry dtos.ExpenseLogs
			nestedRow.Find("td").Each(func(l int, cell *goquery.Selection) {
				text := strings.TrimSpace(cell.Text())
				switch l {
				case 0:
					expenseEntry.Claim = text
				case 1:
					expenseEntry.Supplier = text
				case 2:
					expense, err := format.ExpenseToFloat(text)
					if err != nil {
						parseErr = err
						return
					}
					expenseEntry.Cost = expense
				}
			})
			if expenseEntry.Claim == "" || expenseEntry.Cost == 0 || expenseEntry.Supplier == "" {
				return
			}
			expenseLogs = append(expenseLogs, expenseEntry)
		})
		hospitalityExpense.Event.ExpenseLogs = expenseLogs
		hospitalityExpenses = append(hospitalityExpenses, hospitalityExpense)
	}
	return hospitalityExpenses, parseErr
}

func MpContractExpenses(doc *goquery.Document) ([]*dtos.ContractExpense, error) {
	var contractExpenses []*dtos.ContractExpense
	var parseErr error

	doc.Find("tbody tr").Each(func(i int, row *goquery.Selection) {
		contractExpense := &dtos.ContractExpense{}
		row.Find("td").Each(func(j int, cell *goquery.Selection) {
			text := strings.TrimSpace(cell.Text())

			switch j {
			case 0:
				contractExpense.Supplier = format.Supplier(text)
			case 1:
				contractExpense.Description = text
			case 2:
				contractExpense.Date = text
			case 3:
				expense, err := format.ExpenseToFloat(text)
				if err != nil {
					parseErr = err
					break
				}
				contractExpense.Total = expense
			}
		})
		contractExpenses = append(contractExpenses, contractExpense)
	})

	return contractExpenses, parseErr
}

/*
func MpAllExpenses(mp *dtos.MpWithExpenseCategories, b *browser.Browser) dtos.MpExpensesResults {
	var wg sync.WaitGroup

	travelChan := make(chan []*dtos.TravelExpense, 1)
	errChan := make(chan error, 3)

	wg.Add(1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractTravelExpenses",
			Url:                mp.TravelExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		travelExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			if travelExpenses, ok := travelExpenses.([]*dtos.TravelExpense); ok {
				travelChan <- travelExpenses
			} else {
				errChan <- fmt.Errorf("type assertion failed for hospitality expenses")
			}
		}
	}()

	wg.Add(1)
	hospitalityChan := make(chan []*dtos.HospitalityExpense, 1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractHospitalityExpenses",
			Url:                mp.HospitalityExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		hospitalityExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			if hospitalityExpenses, ok := hospitalityExpenses.([]*dtos.HospitalityExpense); ok {
				hospitalityChan <- hospitalityExpenses
			} else {
				errChan <- fmt.Errorf("type assertion failed for hospitality expenses")
			}
		}
	}()

	wg.Add(1)
	contractChan := make(chan []*dtos.ContractExpense, 1)
	go func() {
		defer wg.Done()

		task := dtos.Task{
			Type:               "extractContractExpenses",
			Url:                mp.ContractExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		contractExpenses, err := b.RunTask(task)
		if err != nil {
			errChan <- err
			return
		} else {
			contractExpenses, ok := contractExpenses.([]*dtos.ContractExpense)
			if !ok {
				errChan <- err
			} else {
				contractChan <- contractExpenses
			}
		}
	}()

	go func() {
		wg.Wait()
		close(travelChan)
		close(hospitalityChan)
		close(contractChan)
		close(errChan)
	}()

	result := dtos.MpExpensesResults{}

	for travel := range travelChan {
		result.TravelExpenses = travel
	}
	for hospitality := range hospitalityChan {
		result.HospitalityExpenses = hospitality
	}
	for contract := range contractChan {
		result.ContractExpenses = contract
	}
	for err := range errChan {
		result.Errors = append(result.Errors, err)
	}
	return result
}*/
