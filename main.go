package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"

	//"log"

	"time"
)

func main() {
	startTime := time.Now()
	browser := &browser.Browser{}

	task := constants.PollingTask
	reports, err := browser.RunTask(task)
	if err != nil {
		fmt.Println("Err")
	}

	expenditureReports := reports.(dtos.AllExpenditureReports)
	for _, report := range expenditureReports.Reports {
		mpTask := constants.MpTask
		mpTask.Url = report.Href

		mps, err := browser.RunTask(mpTask)
		if err != nil {
			fmt.Println("Error extract MP data: ", err)
		}

		mpData := mps.([]*dtos.MpWithExpenseCategories)
		for _, mp := range mpData {

			var travelExpenses []*dtos.TravelExpense
			var contractExpenses []*dtos.ContractExpense
			var hospitalityExpenses []*dtos.HospitalityExpense

			if mp.TravelExpenses.Href != "" {
				tt := constants.TravelTask
				tt.Url = mp.TravelExpenses.Href

				travelData, err := browser.RunTask(tt)
				if err != nil {
					fmt.Println("Travel extraction error: ", err)
				}

				travelExpenses = travelData.([]*dtos.TravelExpense)

			}

			if mp.ContractExpenses.Href != "" {
				ct := constants.ContractTask
				ct.Url = mp.ContractExpenses.Href

				contractData, err := browser.RunTask(ct)
				if err != nil {
					fmt.Println("Contract extraction error: ", err)
				}
				contractExpenses = contractData.([]*dtos.ContractExpense)
			}

			if mp.HospitalityExpenses.Href != "" {
				ht := constants.HospitalityTask
				ht.Url = mp.HospitalityExpenses.Href

				hospitalityData, err := browser.RunTask(ht)
				if err != nil {
					fmt.Println("Hospitality extraction error: ", err)
				}

				hospitalityExpenses = hospitalityData.([]*dtos.HospitalityExpense)
			}

			fmt.Println("Travel expenses: ", travelExpenses)
			fmt.Println("Contract expenses: ", contractExpenses)
			fmt.Println("Hospitality expenses: ", hospitalityExpenses)
		}
	}

	runTime := time.Now().Sub(startTime)
	fmt.Println("Total Runtime: ", runTime)
}
