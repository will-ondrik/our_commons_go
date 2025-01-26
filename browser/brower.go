package browser

import (
	"context"
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"etl_our_commons/extract"
	"fmt"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

type Browser struct {
	allocCtx context.Context
}

func (b *Browser) RunTask(task dtos.Task) (interface{}, error) {
	fmt.Printf("\nRunning task: %s...\n", task.Type)
	fmt.Printf("\nVisiting URL: %s\n", task.Url)
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), constants.CHROME_OPTIONS...)
	defer cancel()

	ctx, cancelCtx := chromedp.NewContext(allocCtx)
	defer cancelCtx()

	doc, err := b.GetHtml(ctx, task)
	if err != nil {
		return nil, fmt.Errorf("Failed to extract html: %v", err)
	}

	formattedData, err := b.GetData(task.Type, doc)
	if err != nil {
		return nil, fmt.Errorf("Failed to format data for task: %s", task.Type)
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
