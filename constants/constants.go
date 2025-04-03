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

	// Flight keywords
	FLIGHT_KEYWORDS = []string{
		"travel",
		"unite the family",
		"to attend a national caucus meeting",
		"to attend a regional or provincial caucus meeting",
		"to attend meetings with stakeholders about business of the House",
		"to attend a conference",
		"to attend training",
		"travel to/from constituency and Ottawa",
	}

	// Fallback for cities without an airport (airport in close proximity)
	AIRPORT_FALLBACK = map[string]string{
		"Eskasoni": "Sydney",
		"Burnaby": "Vancouver",
		"Boundary Bay" : "Vancouver",
		"Duncan": "Victoria",  
		"Brantford": "Hamilton", 
		"Midland": "Toronto",   
		"Fort McMurray": "Edmonton", 
		"Two Hills": "Edmonton", 
		"Vegreville": "Edmonton", 
		"Thompson": "Winnipeg",   
		"Collingwood": "Toronto", 
		"New Glasgow": "Halifax",
		"Penticton": "Kelowna", 
	}

	// IATA code replacements for problematic airport codes
	IATA_REPLACEMENTS = map[string]string{
		"ZBD": "YVR", // Replace Boundary Bay with Vancouver International
		"DUQ": "YYJ", // Replace Duncan with Victoria
		"YFD": "YHM", // Replace Brantford with Hamilton
		"YEE": "YYZ", // Replace Midland/Huronia with Toronto
		"NML": "YEG", // Replace Fort McMurray/Mildred Lake with Edmonton
		"ZSP": "YEG", // Replace St. Paul with Edmonton
		"YGD": "YYZ", // Replace Goderich with Toronto
		"YJM": "YXS", // Replace Fort St. James with Prince George
		"XCM": "YHZ", // Replace Chatham Kent with Halifax for New Glasgow
	}
)

const (
	// Expense categories
	TRAVEL      = "Travel"
	HOSPITALITY = "Hospitality"
	CONTRACT    = "Contract"

	// Transporation modes
	GROUND_TRANSPORTATION   = "Car"
	AIR_TRANSPORTATION = "Plane"

	// Distance limit (if within this range, it cannot be a flight)
	KM_THRESHHOLD = 100

	// Emissions Units
	// Always in Kilograms
	KILOGRAMS = "kg"

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
