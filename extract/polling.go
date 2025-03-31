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
				if len(text) > 9 {
					quarter, err := FiscalQuarter(text)
					if err != nil {
						fmt.Println(err)
					}
					dateRange, err := ReportDates(text)
					if err != nil {
						panic(err)
					}

					fmt.Println("Date range:", dateRange)
					year, err := ReportYear(dateRange)
					if err != nil {
						fmt.Println("Error converting date range to year")
					}

					cell.Find("a").Each(func(j int, link *goquery.Selection) {
						href, exists := link.Attr("href")
						if !exists {
							fmt.Println("No href")
						}
						hrefLink := fmt.Sprintf("%s%s", constants.BASE_URL, href)

						expenditureReport = dtos.ExpenditureReport{
							FiscalYear:     year,
							DateRange: dateRange,
							FiscalQuarter:   quarter,
							Href:      hrefLink,
						}
					})

					expenditureReports.Reports = append(expenditureReports.Reports, expenditureReport)
				}

			}
		})
		for _, report := range expenditureReports.Reports {
			fmt.Printf("Report: %+v\n\n", report)
		}

	}

	return expenditureReports
}
