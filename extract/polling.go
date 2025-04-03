package extract

import (
	"etl_our_commons/constants"
	"etl_our_commons/dtos"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ExpenditureReports(doc *goquery.Document) dtos.AllExpenditureReports {
	var expenditureReports dtos.AllExpenditureReports

	ul := doc.Find("ul")
	if !ul.HasClass("nav navbar-nav") && ul.Has("ul#ce-hoc-nav-parliamentary-business").Length() == 0 {
		ul.Find("li").Each(func(i int, cell *goquery.Selection) {
			var expenditureReport dtos.ExpenditureReport
			if i > 6 {
				text := strings.TrimSpace(cell.Text())
				fmt.Println("\n--- Processing list item ---")
				fmt.Println("Raw text:", text)
				
				if len(text) > 9 {
					quarter, err := FiscalQuarter(text)
					if err != nil {
						fmt.Println("FiscalQuarter error:", err)
					}
					fmt.Println("Fiscal quarter:", quarter)
					
					dateRange, err := ReportDates(text)
					if err != nil {
						fmt.Println("ReportDates error:", err)
						panic(err)
					}

					fmt.Println("Date range after ReportDates:", dateRange)
					
					year, err := ReportYear(dateRange)
					if err != nil {
						fmt.Println("ReportYear error:", err)
					}
					fmt.Println("Year after ReportYear:", year)

					cell.Find("a").Each(func(j int, link *goquery.Selection) {
						href, exists := link.Attr("href")
						if !exists {
							fmt.Println("No href found")
						}
						hrefLink := fmt.Sprintf("%s%s", constants.BASE_URL, href)
						fmt.Println("Link:", hrefLink)

						expenditureReport = dtos.ExpenditureReport{
							FiscalYear:     year,
							DateRange:      dateRange,
							FiscalQuarter:  quarter,
							Href:           hrefLink,
						}
						
						fmt.Println("Created expenditure report:", expenditureReport)
					})

					expenditureReports.Reports = append(expenditureReports.Reports, expenditureReport)
				}
			}
		})
		
		fmt.Println("\n--- All Reports ---")
		for i, report := range expenditureReports.Reports {
			fmt.Printf("Report %d: %+v\n\n", i, report)
		}
	}

	return expenditureReports
}
