package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"etl_our_commons/tasks"
	"fmt"
	"sync"
	"time"
)

const (
	WorkerLimit       = 2   // Number of concurrent workers
	RequestsPerSecond = 2   // Rate limit for external requests
	mpQueueBuffer     = 100 // Increased channel buffer to decouple production from consumption
)

func main() {
	startTime := time.Now()
	b := &browser.Browser{}

	tm := tasks.NewTaskManager(b)

	expenditures, err := tm.PollForReports()
	if err != nil {
		panic(err)
	}

	// Create a buffered channel for MPs and set up a ticker-based rate limiter.
	mpQueue := make(chan *dtos.MpWithExpenseCategories, mpQueueBuffer)
	ticker := time.NewTicker(time.Second / RequestsPerSecond)
	defer ticker.Stop()
	rateLimiter := ticker.C

	var wg sync.WaitGroup

	// Start the worker pool.
	for i := 0; i < WorkerLimit; i++ {
		wg.Add(1)
		go tm.ProcessMpQueue(mpQueue, &wg, rateLimiter)
	}

	// Enqueue MPs from the expenditure reports (processing only the first report as before).
	for i, report := range expenditures.Reports {
		if i == 1 {
			break
		}
		mps, err := tm.MpExpenditures(report.Href)
		if err != nil {
			panic(err)
		}

		for _, mp := range mps {
			mpQueue <- mp
		}
	}

	close(mpQueue)
	wg.Wait()

	runTime := time.Since(startTime)
	fmt.Println("Total Runtime:", runTime)
}
