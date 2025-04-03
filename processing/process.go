package processing

import (
	"etl_our_commons/dtos"
	"etl_our_commons/extract"
	"fmt"
	"log"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

func GetData(taskType string, doc *goquery.Document) (interface{}, error) {
	var output interface{}
	var taskErr error

	switch taskType {

	case "extractMps":
		formattedMpData, err := extract.Mps(doc)
		if err != nil {
			taskErr = err
		}
		output = formattedMpData

	case "extractTravelExpenses":
		formattedTravelData, err := extract.MpTravelExpenses(doc)
		if err != nil {
			taskErr = err
		}
		output = formattedTravelData
	case "extractContractExpenses":
		formattedContractData, err := extract.MpContractExpenses(doc)
		if err != nil {
			taskErr = err
		}
		output = formattedContractData
	case "extractHospitalityExpenses":
		formattedHospitalityData, err := extract.MpHospitalityExpenses(doc)
		if err != nil {
			taskErr = err
		}
		output = formattedHospitalityData
	case "polling":
		formattedReports := extract.ExpenditureReports(doc)
		output = formattedReports

	default:
		log.Printf("Unknown task: %s", taskType)

	}

	if taskErr != nil {
		return nil, fmt.Errorf("Failed to extract task data: %v", taskErr)
	}

	return output, nil
}

func ProcessData(mpHtml []*dtos.MpHtml) []*dtos.Mp {
	var mpsExtracted []*dtos.Mp
	var wg sync.WaitGroup

	for _, html := range mpHtml {
		fmt.Printf("HTML: %+v\n", html)
		mp := &dtos.Mp{}

		mp.MpName = html.Info.MpName
		mp.Caucus = html.Info.Caucus
		mp.Constituency = html.Info.Constituency
		mp.Expenses.Totals.SalariesCost = html.Info.Salaries
		mp.Expenses.Totals.ContractCost = html.Info.ContractExpenses.ExpenseTotal
		mp.Expenses.Totals.HospitalityCost = html.Info.HospitalityExpenses.ExpenseTotal
		mp.Expenses.Totals.TravelCost = html.Info.TravelExpenses.ExpenseTotal
		mp.Years = html.Info.Years
		mp.FiscalYear = html.Info.FiscalYear
		mp.FiscalQuarter = html.Info.FiscalQuarter
		mp.Url = html.Info.Url 
		
		// Add expense URLs to the Expenses struct
		mp.Expenses.ContractExpensesUrl = html.Info.ContractExpenses.Href
		mp.Expenses.HospitalityExpensesUrl = html.Info.HospitalityExpenses.Href
		mp.Expenses.TravelExpensesUrl = html.Info.TravelExpenses.Href

		if html.Contract != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := extract.MpContractExpenses(html.Contract)
				if err != nil {
					panic(err)
				}
				mp.Expenses.ContractExpenses = data
			}()
		}

		if html.Hospitality != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := extract.MpHospitalityExpenses(html.Hospitality)
				if err != nil {
					panic(err)
				}
				mp.Expenses.HospitalityExpenses = data
			}()

		}

		if html.Travel != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				data, err := extract.MpTravelExpenses(html.Travel)
				if err != nil {
					panic(err)
				}
				mp.Expenses.TravelExpenses = data
			}()
		}

		wg.Wait()
		mpsExtracted = append(mpsExtracted, mp)
	}

	return mpsExtracted
}
