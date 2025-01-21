package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"fmt"
	"log"
)

func main() {
	browser := &browser.Browser{}

	task := dtos.Task{
		Type:               "extractMps",
		Url:                "https://www.ourcommons.ca/proactivedisclosure/en/members/2022/1",
		ExtractFromElement: "table.table",
	}

	mpData, err := browser.RunTask(task)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(mpData.([]*dtos.MpWithExpenseCategories))
	for _, mp := range mpData.([]*dtos.MpWithExpenseCategories) {
		fmt.Printf("\nMP: %+v\n", mp)
	}

	for i, mp := range mpData.([]*dtos.MpWithExpenseCategories) {
		if i == 0 {
			task = dtos.Task{
				Type:               "extractTravelExpenses",
				Url:                mp.TravelExpenses.Href,
				ExtractFromElement: "#data-table",
			}

			mpTravelExpenses, err := browser.RunTask(task)
			if err != nil {
				fmt.Printf("\nTravel expenses extraction failed: %v", err)
			}

			fmt.Println(mpTravelExpenses.([]*dtos.TravelExpenses))
		}

	}

}
