package browser

import (
	"context"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"etl_our_commons/extract"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Browser struct {
	allocCtx context.Context
}

func (b *Browser) RunTask(task dtos.Task) (interface{}, error) {
	fmt.Println("[RUNNING TASK]: ", task.Type)
	fmt.Println("[VISITING URL]: ", task.Url)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), constants.CHROME_OPTIONS...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	//doc, err := b.GetHtml(ctx, task)
	doc, err := b.GetHtml2(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract html: %v\n", err)
	}

	formattedData, err := b.GetData(task.Type, doc)
	if err != nil {
		return nil, fmt.Errorf("Failed to format data for task: %s\n", task.Type)
	}

	return formattedData, nil
}

func (b *Browser) GetData(taskType string, doc *goquery.Document) (interface{}, error) {
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

// TODO: Need to handle timeouts or error page
// Re-run function if either occur
func (b *Browser) GetHtml(ctx context.Context, task dtos.Task) (*goquery.Document, error) {
	var html string
	fmt.Println("Looking for html element...")

	err := chromedp.Run(ctx,
		chromedp.Navigate(task.Url),
		chromedp.WaitVisible(task.ExtractFromElement),
		chromedp.OuterHTML(task.ExtractFromElement, &html),
	)

	if err != nil {
		return nil, fmt.Errorf("Chrome instance failed: %v", err)
	}

	if html != "" {
		fmt.Println("html found")
	}
	err = chromedp.Cancel(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to close Chrome instance: %v", err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("Unable to parse html: %v", err)
	}
	fmt.Println("Returning doc...")
	return doc, nil
}

func (b *Browser) GetHtml2(ctx context.Context, task dtos.Task) (*goquery.Document, error) {
	var html string

	attempts := 0
	maxAttempts := 5

	for attempts < maxAttempts {
		// Create a new child context for this attempt.
		childCtx, cancel := chromedp.NewContext(ctx)
		fmt.Printf("[ATTEMPT: %d]...\n", attempts)
		attempts++

		var err error
		if task.ExtractFromElement == "body" {
			err = chromedp.Run(childCtx,
				chromedp.Navigate(task.Url),
				chromedp.WaitReady(task.ExtractFromElement),
				chromedp.OuterHTML(task.ExtractFromElement, &html),
			)
		} else {
			err = chromedp.Run(childCtx,
				chromedp.Navigate(task.Url),
				chromedp.WaitVisible(task.ExtractFromElement),
				chromedp.OuterHTML(task.ExtractFromElement, &html),
			)
		}

		// Cancel the child context immediately.
		cancel()

		if err != nil {
			fmt.Println("[CHROME INSTANCE FAILED] Retrying...")
			continue
		}

		if html == "" {
			fmt.Println("[EMPTY HTML] Retrying...")
			continue
		}

		if b.IsErrorPage(html) {
			fmt.Println("[RATE LIMITED] Sleeping...")
			time.Sleep(2 * time.Second)
			fmt.Println("[SLEEP OVER] Retrying...")
			continue
		}

		if html != "" {
			fmt.Println("[HTML FOUND] Done!")
			break
		}
	}

	if html == "" {
		fmt.Println("[FAILED TASK]")
		return nil, fmt.Errorf("all attempts exhausted. task failed.\n")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	fmt.Println("Returning doc...")
	return doc, nil
}

func (b *Browser) IsErrorPage(html string) bool {
	return strings.Contains(html, "System Error")
}

func (b *Browser) CancelInstance(ctx context.Context) {
	chromedp.Cancel(ctx)
	return
}
