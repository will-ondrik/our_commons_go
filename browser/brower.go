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
	fmt.Printf("[RUNNING TASK]: %s\n", task.Type)
	fmt.Printf("[VISTING URL]: %s\n", task.Url)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), constants.CHROME_OPTIONS...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	doc, err := b.GetHtml(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract html: %v\n", err)
	}

	formattedData, err := b.GetData(task.Type, doc)
	if err != nil {
		return nil, fmt.Errorf("Failed to format data for task: %s\n", task.Type)
	}

	return formattedData, nil
}

func (b *Browser) RunTask2(task dtos.Task) (*goquery.Document, error) {
	fmt.Printf("[RUNNING TASK]: %s\n", task.Type)
	fmt.Printf("[VISTING URL]: %s\n", task.Url)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), constants.CHROME_OPTIONS...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	if doc, err := b.GetHtml(ctx, task); err != nil {
		return nil, fmt.Errorf("Failed to extract html: %v\n", err)
	} else {
		return doc, nil
	}
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

func (b *Browser) GetHtml(ctx context.Context, task dtos.Task) (*goquery.Document, error) {
	var html string

	attempts := 0
	max_attempts := 5

	for attempts < max_attempts {
		ctx, cancel := chromedp.NewContext(ctx)
		defer cancel()
		fmt.Printf("\n[ATTEMPT: %d]...\n", attempts)
		attempts++

		var err error
		if task.ExtractFromElement == "body" {
			err = chromedp.Run(ctx,
				chromedp.Navigate(task.Url),
				chromedp.WaitReady(task.ExtractFromElement),
				chromedp.OuterHTML(task.ExtractFromElement, &html),
			)
		} else {
			err = chromedp.Run(ctx,
				chromedp.Navigate(task.Url),
				chromedp.WaitVisible(task.ExtractFromElement),
				chromedp.OuterHTML(task.ExtractFromElement, &html),
			)
		}

		if err != nil {
			fmt.Println("[CHROME INSTANCE FAILED] Retrying...")
			b.CancelInstance(ctx)
			continue
		}

		if html == "" {
			fmt.Println("[EMPTY HTML] Retrying...")
			b.CancelInstance(ctx)
			continue
		}
		if b.IsErrorPage(html) {
			fmt.Println("[RATE LIMITED] Sleeping...")
			b.CancelInstance(ctx)
			time.Sleep(2 * time.Second)
			fmt.Println("[SLEEP OVER] Retrying...")
			continue
		}

		if html != "" {
			fmt.Println("[HTML FOUND] Done!")
			b.CancelInstance(ctx)
			break
		}

	}

	if html == "" || attempts == 4 {
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
	if strings.Contains(html, "System Error") {
		return true
	}

	return false
}

func (b *Browser) CancelInstance(ctx context.Context) {
	chromedp.Cancel(ctx)
	return
}
