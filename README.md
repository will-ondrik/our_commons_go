# Our Commons Go

## Overview

**Our Commons Go** is a Go-based web scraper that extracts **Member of Parliament (MP) expenditures** from publicly available government sources. This project aims to provide structured data on how Canadian MPs allocate their spending.

The scraper is built using:
- **ChromeDP** (Headless Chrome) for navigating government sites.
- **GoQuery** for parsing and extracting relevant HTML elements.

Future updates will expand the project to include:
- MP voting records
- MP roles and committee positions
- House officers and administration details

## Features

✔ **Scrapes MP expenditures** from official sources.  
✔ **Formats data into a structured output** for further analysis.  
✔ **Handles dynamic content** using headless Chrome.  
✔ **Optimized for automation**, allowing scheduled scraping.

## Installation

### Prerequisites
- **Go 1.18+** installed
- **Google Chrome** (for ChromeDP)
- **Chromedriver** (ensure compatibility with your Chrome version)

### Steps

1. **Clone the repository:**
   ```bash
   git clone https://github.com/will-ondrik/our_commons_go.git
   cd our_commons_go
