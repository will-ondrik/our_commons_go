package main

import (
	"context"
	"etl_our_commons/dtos"
	"etl_our_commons/extract"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/chromedp"
)

func main() {
	// Set Chrome options
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false),
		chromedp.Flag("disable-gpu", false),
	)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()

	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()

	ctx, cancel = context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	url := "https://www.ourcommons.ca/proactivedisclosure/en/members"
	var tableHTML string

	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitVisible("table.table"),
		chromedp.OuterHTML("table.table", &tableHTML),
	)
	if err != nil {
		fmt.Println("Error during Chromedp:", err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(tableHTML))

	// extract MP Expense Categories
	var mps []*dtos.MpWithExpenseCategories
	mps, err = extract.Mps(doc)
	if err != nil {
		log.Printf("Error extracting MP data: %v", err)
	}

	fmt.Println(mps)

	cancel()
}
