package tasks

import (
	"etl_our_commons/dtos"
)

var (
	MpTask = dtos.Task{
		Type:               "extractMps",
		Url:                "",
		ExtractFromElement: "body",
	}

	TravelTask = dtos.Task{
		Type:               "extractTravelExpenses",
		Url:                "",
		ExtractFromElement: "body",
	}

	HospitalityTask = dtos.Task{
		Type:               "extractHospitalityExpenses",
		Url:                "",
		ExtractFromElement: "body",
	}

	ContractTask = dtos.Task{
		Type:               "extractContractExpenses",
		Url:                "",
		ExtractFromElement: "body",
	}

	PollingTask = dtos.Task{
		Type:               "polling",
		Url:                "https://www.ourcommons.ca/proactivedisclosure/en/members/2022/1",
		ExtractFromElement: "main.ce-hoc-body-content",
	}
)
