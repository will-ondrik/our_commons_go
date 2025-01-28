package constants

import (
	"etl_our_commons/dtos"

	"github.com/chromedp/chromedp"
)

var (

	// Expense categories
	TRAVEL      = "Travel"
	HOSPITALITY = "Hospitality"
	CONTRACT    = "Contract"

	// Base MP expenditures URL
	BASE_URL = "https://www.ourcommons.ca"

	// Naming
	VACANT = "Vacant"

	CHROME_OPTIONS = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
	)

	MpTask = dtos.Task{
		Type:               "extractMps",
		Url:                "",
		ExtractFromElement: "#data-table",
	}

	TravelTask = dtos.Task{
		Type:               "extractTravelExpenses",
		Url:                "",
		ExtractFromElement: "#data-table",
	}

	HospitalityTask = dtos.Task{
		Type:               "extractHospitalityExpenses",
		Url:                "",
		ExtractFromElement: "#data-table",
	}

	ContractTask = dtos.Task{
		Type:               "extractContractExpenses",
		Url:                "",
		ExtractFromElement: "#data-table",
	}

	PollingTask = dtos.Task{
		Type:               "polling",
		Url:                "https://www.ourcommons.ca/proactivedisclosure/en/members/2022/1",
		ExtractFromElement: "main.ce-hoc-body-content",
	}
)
