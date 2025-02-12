package tasks

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"fmt"
	"log"
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

func (tm *TaskManager) ProcessMpHtml(tasks dtos.Tasks) (*dtos.MpHtml, error, dtos.TaskErrors) {
	log.Printf("Processing MP html...\n")
	mpHtml := &dtos.MpHtml{}
	var wg sync.WaitGroup

	var taskErrors dtos.TaskErrors
	errChan := make(chan error, 3)

	if tasks.Contract.Url == "" {
		mpHtml.Contract = nil
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := tm.Browser.RunTask2(tasks.Contract)
			if err != nil {
				fmt.Println("[CONTRACT ERROR]", tasks.Contract.Url)
				taskErrors.Contract = true
				errChan <- err
			} else {
				mpHtml.Contract = doc
			}
		}()
	}

	if tasks.Hospitality.Url == "" {
		mpHtml.Hospitality = nil
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond)
			doc, err := tm.Browser.RunTask2(tasks.Hospitality)
			if err != nil {
				fmt.Println("[HOSPITALITY ERROR]", tasks.Hospitality.Url)
				taskErrors.Hospitality = true
				errChan <- err
			} else {
				mpHtml.Hospitality = doc

			}
		}()
	}

	if tasks.Travel.Url == "" {
		mpHtml.Travel = nil
	} else {
		wg.Add(1)
		go func() {
			defer wg.Done()
			time.Sleep(200 * time.Millisecond)
			doc, err := tm.Browser.RunTask2(tasks.Travel)
			if err != nil {
				fmt.Println("[TRAVEL ERROR]", tasks.Hospitality.Url)
				taskErrors.Travel = true
				errChan <- err
			} else {
				mpHtml.Travel = doc
			}
		}()
	}

	wg.Wait()

	// Collect and close error channel
	close(errChan)
	var err error
	for e := range errChan {
		err = e
	}

	return mpHtml, err, taskErrors
}

func (tm *TaskManager) CreateTasks(mp *dtos.MpWithExpenseCategories) dtos.Tasks {
	ct := ContractTask
	ct.Url = mp.ContractExpenses.Href

	ht := HospitalityTask
	ht.Url = mp.HospitalityExpenses.Href

	tt := TravelTask
	tt.Url = mp.TravelExpenses.Href

	return dtos.Tasks{
		Contract:    ct,
		Hospitality: ht,
		Travel:      tt,
	}
}

func (tm *TaskManager) ProcessMpQueue2(mpQueue chan *dtos.MpWithExpenseCategories, wg *sync.WaitGroup, rateLimiter <-chan time.Time, htmlStore *[]*dtos.MpHtml, mu *sync.Mutex, redoTaskQueue *[]*dtos.RedoTask) {
	defer wg.Done()
	for mp := range mpQueue {
		<-rateLimiter
		log.Printf("Process MPQ2 time for html")
		tasks := tm.CreateTasks(mp)

		mpHtml, err, taskErrs := tm.ProcessMpHtml(tasks)
		mpHtml.Info = *mp

		if err != nil {
			log.Printf("[ERROR] Failed to scrape MP HTML")
			// TODO: Add failed MPs to a slice to revisit
			tm.AddRedoTask(redoTaskQueue, *mpHtml, tasks, taskErrs)

			fmt.Printf("\n[TASK QUEUE] Added task to revisit: %v\n", tasks)
			continue
		}
		log.Printf("[SUCCESS] Item queued")
		mu.Lock()
		*htmlStore = append(*htmlStore, mpHtml)
		mu.Unlock()
	}
}

func (tm *TaskManager) AddRedoTask(redoQueue *[]*dtos.RedoTask, mpHtml dtos.MpHtml, tasks dtos.Tasks, taskErrs dtos.TaskErrors) {
	*redoQueue = append(*redoQueue, &dtos.RedoTask{
		MpHtml:     mpHtml,
		Tasks:      tasks,
		TaskErrors: taskErrs,
	})
}

func (tm *TaskManager) RedoTask(task *dtos.RedoTask) (dtos.MpHtml, error) {
	var wg sync.WaitGroup

	if task.TaskErrors.Contract {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := tm.Browser.RunTask2(task.Tasks.Contract)
			if err != nil {
				fmt.Println("[REDO CONTRACT ERROR]")
			} else {
				task.MpHtml.Contract = doc
			}
		}()
	}

	if task.TaskErrors.Hospitality {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := tm.Browser.RunTask2(task.Tasks.Hospitality)
			if err != nil {
				fmt.Println("[REDO HOSPITALITY ERROR]")
			} else {
				task.MpHtml.Hospitality = doc
			}
		}()
	}

	if task.TaskErrors.Travel {
		wg.Add(1)
		go func() {
			defer wg.Done()
			doc, err := tm.Browser.RunTask2(task.Tasks.Travel)
			if err != nil {
				fmt.Println("[REDO TRAVEL ERROR]")
			} else {
				task.MpHtml.Travel = doc
			}
		}()
	}
	wg.Wait()

	return task.MpHtml, nil
}

func (tm *TaskManager) ProcessRedoTasks(redoTasks []*dtos.RedoTask, htmlStore *[]*dtos.MpHtml) {
	for _, task := range redoTasks {
		mpHtml, err := tm.RedoTask(task)
		if err != nil {
			fmt.Println(err)
		}
		*htmlStore = append(*htmlStore, &mpHtml)
	}
}
