package main

import (
	"etl_our_commons/browser"
	"etl_our_commons/dtos"
	"fmt"
	"log"
	"time"
)

func main() {
	startTime := time.Now()
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

	for _, mp := range mpData.([]*dtos.MpWithExpenseCategories) {

		if mp.HospitalityExpenses.Href == "nil" {
			fmt.Println("MP is nil: ", mp)
			continue
		}
		// TODO: // If href is nil, skip task
		task = dtos.Task{
			Type:               "extractHospitalityExpenses",
			Url:                mp.HospitalityExpenses.Href,
			ExtractFromElement: "#data-table",
		}

		hospitalityExpenses, err := browser.RunTask(task)
		if err != nil {
			fmt.Printf("\nHospitality expenses extraction failed: %v", err)
		}

		for _, te := range hospitalityExpenses.([]*dtos.HospitalityExpense) {
			fmt.Printf("\n%+v\n", te)
		}

	}

	runTime := time.Now().Sub(startTime)
	fmt.Println("Total Runtime: ", runTime)

}
