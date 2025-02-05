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
		liSelection := ul.Find("li")
		// Preallocate based on number of li elements.
		reports := make([]dtos.ExpenditureReport, 0, liSelection.Length())
		liSelection.Each(func(i int, cell *goquery.Selection) {
			var expenditureReport dtos.ExpenditureReport
			if i > 6 {
				text := strings.TrimSpace(cell.Text())
				if len(text) > 9 {
					quarter, err := FiscalQuarter(text)
					if err != nil {
						fmt.Println(err)
					}
					dateRange := ReportDates(text)
					yearRange := ReportYears(dateRange)

					cell.Find("a").Each(func(j int, link *goquery.Selection) {
						href, exists := link.Attr("href")
						if !exists {
							fmt.Println("No href")
						}
						hrefLink := fmt.Sprintf("%s%s", constants.BASE_URL, href)
						expenditureReport = dtos.ExpenditureReport{
							Years:     yearRange,
							DateRange: dateRange,
							Quarter:   quarter,
							Href:      hrefLink,
						}
					})

					reports = append(reports, expenditureReport)
				}
			}
		})
		expenditureReports.Reports = reports

		for _, report := range expenditureReports.Reports {
			fmt.Printf("Report: %+v\n\n", report)
		}
	}

	return expenditureReports
}
