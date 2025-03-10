package constants

import (
	"github.com/chromedp/chromedp"
)

var (

	CHROME_OPTIONS = append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
	)

	// Scraping attempt limits
	ATTEMPTS = 0
	MAX_ATTEMPTS = 5

	// Urls for MP contact information
	PARLIAMENT_MEMBERS_INFO = map[string]string{
		"PARLIAMENT_43": "https://www.ourcommons.ca/Members/en/search?parliament=43&caucusId=all&province=all&gender=all",
		"PARLIAMENT_44": "https://www.ourcommons.ca/Members/en/search?parliament=44&caucusId=all&province=all&gender=all",
	}
)

const (
	// Expense categories
	TRAVEL      = "Travel"
	HOSPITALITY = "Hospitality"
	CONTRACT    = "Contract"

	// Base URL
	BASE_URL = "https://www.ourcommons.ca"

	// Naming
	VACANT = "Vacant"

	// HTML scraping keywords
	HTML_BODY = "body"
	EMPTY_STR = ""

	// Worker and rate limits
	// for Chrome scraping
	WORKER_LIMIT = 4
	REQUESTS_PER_SECOND = 2


)
