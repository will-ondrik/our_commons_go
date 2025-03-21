package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"etl_our_commons/processing"
	"etl_our_commons/tasks"
	"fmt"
	"sync"
	"time"
)

/*
Runtimes
- Sequential runtime (commented out): 47 minutes per report
- Concurrent runtime (updated code): 24 minutes, 36 seconds per report
*/

// Worker and Rate Limits
const WorkerLimit = constants.WORKER_LIMIT
const RequestsPerSecond = constants.REQUESTS_PER_SECOND

func main() {
	startTime := time.Now()
	defer getRuntime(startTime)
	b := &browser.Browser{}

	tm := tasks.NewTaskManager(b)

	expenditures, err := tm.PollForReports()
	if err != nil {
		panic(err)
	}

	// TODO: Only process new reports
	// Check DB for list of reports before processing


	var htmlStore []*dtos.MpHtml
	var redoTaskQueue []*dtos.RedoTask
	var mu sync.Mutex

	// Create channel and rate limiter
	// Rate limited ensures set time between requests
	mpQueue := make(chan *dtos.MpWithExpenseCategories, WorkerLimit)
	rateLimiter := time.Tick(time.Second / RequestsPerSecond)
	var wg sync.WaitGroup

	// Create worker pool
	// TODO: Test runtime with larger pool
	for i := 0; i < WorkerLimit; i++ {
		wg.Add(1)
		//go tm.ProcessMpQueue(mpQueue, &wg, rateLimiter)
		go tm.ProcessMpQueue2(mpQueue, &wg, rateLimiter, &htmlStore, &mu, &redoTaskQueue)
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
			mp.Year = report.Years
			mp.Quarter = report.Quarter
			mpQueue <- mp
		}
	}

	// Once all MPs are enqueued, close pool
	// Wait for all processing to complete
	close(mpQueue)
	wg.Wait()

	// Reprocess failed tasks
	fmt.Println("\nHandling failed tasks...")
	tm.ProcessRedoTasks(redoTaskQueue, &htmlStore)

	// Process extracted HTML for Mps
	fmt.Println("\nProcessing html store...")
	mps := processing.ProcessData(htmlStore)
	for _, mp := range mps {
		fmt.Println("-----------------------------------")
		fmt.Printf("MP Info: %+v\n", mp)
		fmt.Println("-----------------------------------")


		// TODO: Save items to db

		
	}
}

func getRuntime(elapsed time.Time) {
	fmt.Println("Total Runtime: ", time.Since(elapsed))
}
