package tasks

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"fmt"
	"sync"
	"time"
)

type TaskManager struct {
	Browser *browser.Browser
}

func NewTaskManager(b *browser.Browser) *TaskManager {
	return &TaskManager{
		Browser: b,
	}
}

// Extract Contract Expenses
func (tm *TaskManager) ContractExpenses(url string) ([]*dtos.ContractExpense, error) {
	t := ContractTask
	t.Url = url
	data, err := tm.Browser.RunTask(t)
	if err != nil {
		return nil, err
	}

	expenses, ok := data.([]*dtos.ContractExpense)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for contract expenses")
	}
	return expenses, nil
}

// Extract Hospitality Expenses
func (tm *TaskManager) HospitalityExpenses(url string) ([]*dtos.HospitalityExpense, error) {
	t := HospitalityTask
	t.Url = url
	data, err := tm.Browser.RunTask(t)
	if err != nil {
		return nil, err
	}

	expenses, ok := data.([]*dtos.HospitalityExpense)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for hospitality expenses")
	}
	return expenses, nil
}

// Extract Travel Expenses
func (tm *TaskManager) TravelExpenses(url string) ([]*dtos.TravelExpense, error) {
	t := TravelTask
	t.Url = url
	data, err := tm.Browser.RunTask(t)
	if err != nil {
		return nil, err
	}

	expenses, ok := data.([]*dtos.TravelExpense)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for travel expenses")
	}
	return expenses, nil
}

// Extract MP Expenditures
func (tm *TaskManager) MpExpenditures(url string) ([]*dtos.MpWithExpenseCategories, error) {
	t := MpTask
	t.Url = url
	data, err := tm.Browser.RunTask(t)
	if err != nil {
		return nil, err
	}

	mpExpenditures, ok := data.([]*dtos.MpWithExpenseCategories)
	if !ok {
		return nil, fmt.Errorf("type assertion failed for MP expenditures")
	}
	return mpExpenditures, nil
}

// Poll for Reports
// TODO: Add func to server
// Compare report list to processed reports
// Process any new reports
func (tm *TaskManager) PollForReports() (dtos.AllExpenditureReports, error) {
	var reports dtos.AllExpenditureReports
	t := PollingTask
	data, err := tm.Browser.RunTask(t)
	if err != nil {
		return reports, err
	}

	reports, ok := data.(dtos.AllExpenditureReports)
	if !ok {
		return reports, fmt.Errorf("type assertion failed for expenditure reports")
	}
	return reports, nil
}

// Process MP's expenses concurrently
// TODO: Return an empty slice if no URL exists
// TODO: Return slice
// TODO: Send returned slice to Kafka for processing
func (tm *TaskManager) ProcessMp(mp *dtos.MpWithExpenseCategories) {
	var wg sync.WaitGroup

	if mp.ContractExpenses.Href != "" {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_, err := tm.ContractExpenses(url)
			if err != nil {
				fmt.Println("[ERROR] Contract Expenses:", err)
			} else {
				fmt.Println("[SUCCESS] Contract Expenses Extracted")
			}
		}(mp.ContractExpenses.Href)
	}

	if mp.HospitalityExpenses.Href != "" {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_, err := tm.HospitalityExpenses(url)
			if err != nil {
				fmt.Println("[ERROR] Hospitality Expenses:", err)
			} else {
				fmt.Println("[SUCCESS] Hospitality Expenses Extracted")
			}
		}(mp.HospitalityExpenses.Href)
	}

	if mp.TravelExpenses.Href != "" {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			_, err := tm.TravelExpenses(url)
			if err != nil {
				fmt.Println("[ERROR] Travel Expenses:", err)
			} else {
				fmt.Println("[SUCCESS] Travel Expenses Extracted")
			}
		}(mp.TravelExpenses.Href)
	}

	wg.Wait()
}

// Worker pool to process MPs
// TODO: Test rate limited and worker pool
// Check if its possible to increase size
func (tm *TaskManager) ProcessMpQueue(mpQueue chan *dtos.MpWithExpenseCategories, wg *sync.WaitGroup, rateLimiter <-chan time.Time) {
	for mp := range mpQueue {
		<-rateLimiter
		tm.ProcessMp(mp)
	}
	wg.Done()
}
