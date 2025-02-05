package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"etl_our_commons/tasks"
	"fmt"
	"sync"
	"time"
)

/*
func main() {
	startTime := time.Now()
	browser := &browser.Browser{}

	task := constants.PollingTask
	reports, err := browser.RunTask(task)
	if err != nil {
		fmt.Println("Err")
	}

	expenditureReports := reports.(dtos.AllExpenditureReports)
	for i, report := range expenditureReports.Reports {
		if i == 1 {
			break
		}
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
*/

/*
Runtimes
- Sequential runtime (commented out): 47 minutes
- Concurrent runtime (updated code): 27 minutes, 36 seconds
*/

// Worker and Rate Limits
const WorkerLimit = 2
const RequestsPerSecond = 2

func main() {
	startTime := time.Now()
	b := &browser.Browser{}

	tm := tasks.NewTaskManager(b)

	expenditures, err := tm.PollForReports()
	if err != nil {
		panic(err)
	}

	// Create channel and rate limiter
	// Rate limited ensures set time between requests
	mpQueue := make(chan *dtos.MpWithExpenseCategories, WorkerLimit)
	rateLimiter := time.Tick(time.Second / RequestsPerSecond)
	var wg sync.WaitGroup

	// Create worker pool
	// TODO: Test runtime with larger pool
	for i := 0; i < WorkerLimit; i++ {
		wg.Add(1)
		go tm.ProcessMpQueue(mpQueue, &wg, rateLimiter)
	}

	// Extract and add MPs to the processing queue
	for i, report := range expenditures.Reports {
		if i == 1 {
			break
		}
		mps, err := tm.MpExpenditures(report.Href)
		if err != nil {
			panic(err)
		}

		for _, mp := range mps {
			// Send MPs to pool
			mpQueue <- mp
		}
	}

	// Once all MPs are enqueued, close pool
	// Wait for all processing to complete
	close(mpQueue)
	wg.Wait()

	runTime := time.Since(startTime)
	fmt.Println("Total Runtime:", runTime)
}
