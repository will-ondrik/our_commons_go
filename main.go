package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/constants"
	"etl_our_commons/database"
	"etl_our_commons/dtos"
	flight "etl_our_commons/flights"
	"etl_our_commons/processing"
	"etl_our_commons/tasks"
	"fmt"
	"runtime"
	"sync"
	"time"
)

/*
Runtimes
- Sequential runtime: 47 minutes per report
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

	// Process each report completely before moving to the next
	for _, report := range expenditures.Reports {
		
		fmt.Printf("\n==== Processing Report: %s ====\n", report.DateRange)
		
		// Create new data stores for this report
		var htmlStore []*dtos.MpHtml
		var redoTaskQueue []*dtos.RedoTask
		var mu sync.Mutex

		// Create channel and rate limiter for this report
		mpQueue := make(chan *dtos.MpWithExpenseCategories, WorkerLimit)
		rateLimiter := time.Tick(time.Second / RequestsPerSecond)
		var wg sync.WaitGroup

		// Create worker pool for this report
		for j := 0; j < WorkerLimit; j++ {
			wg.Add(1)
			go tm.ProcessMpQueue2(mpQueue, &wg, rateLimiter, &htmlStore, &mu, &redoTaskQueue)
		}

		// Extract MPs for this report
		fmt.Println("Extracting MPs for report:", report.DateRange)
		mps, err := tm.MpExpenditures(report.Href)
		if err != nil {
			fmt.Printf("Error extracting MPs for report %s: %v\n", report.DateRange, err)
			continue // Skip to next report on error
		}

		// Queue MPs for processing
		for _, mp := range mps {
			// Send MPs to pool
			mp.Years = report.DateRange
			mp.FiscalQuarter = report.FiscalQuarter
			mp.FiscalYear = report.FiscalYear
			mp.Url = report.Href // Add the report URL to the MP
			mpQueue <- mp
		}

		// Close queue and wait for processing to complete
		close(mpQueue)
		wg.Wait()

		// Reprocess failed tasks for this report
		fmt.Println("\nHandling failed tasks for report:", report.DateRange)
		tm.ProcessRedoTasks(redoTaskQueue, &htmlStore)

		// Process extracted HTML for MPs in this report
		fmt.Println("\nProcessing html store for report:", report.DateRange)
		processedMps := processing.ProcessData(htmlStore)
		fmt.Println("Processing complete for report:", report.DateRange)

		// Initialize flight manager for this report
		fmt.Println("Initializing flight manager for report:", report.DateRange)
		fm, err := flight.NewFlightManager()
		if err != nil {
			fmt.Printf("Error initializing flight manager for report %s: %v\n", report.DateRange, err)
			continue
		}
		
		// Fetch and append travel data to MPs in this report
		fmt.Println("Appending travel data to MPs for report:", report.DateRange)
		updatedMps, err := fm.AppendTravelDataToMps(processedMps)
		if err != nil {
			fmt.Printf("Error appending travel data for report %s: %v\n", report.DateRange, err)
			continue
		}

		// Init db connection and insert data for this report
		fmt.Println("Inserting data into database for report:", report.DateRange)
		db := database.NewDb()
		err = db.InsertAll(db.Pool, updatedMps)
		if err != nil {
			fmt.Printf("Error inserting data for report %s: %v\n", report.DateRange, err)
			continue
		}

		// Write flight cache to file
		// Will heavily cut down on future paid API calls
		// Will have the majority of flight routes
		err = fm.FlightMapToJsonFile(fm.FlightService.Cache)
		if err != nil {
			fmt.Printf("Error writing flight cache for report %s: %v\n", report.DateRange, err)
		}

		fmt.Printf("==== Completed processing report: %s ====\n", report.DateRange)
		
		// Clear references to large data structures
		htmlStore = nil
		redoTaskQueue = nil
		processedMps = nil
		updatedMps = nil
		
		// Force garbage collection after each report
		// Memory was being overwhelmed in previous version
		runtime.GC()
	}

	fmt.Println("All reports processed successfully")
}

func getRuntime(elapsed time.Time) {
	fmt.Println("Total Runtime: ", time.Since(elapsed))
}
