package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/dadosjusbr/coletores"
	"github.com/dadosjusbr/coletores/status"
	"github.com/kelseyhightower/envconfig"
)

type Environment struct {
	Month        int
	Year         int
	OutputFolder string
	GitCommit    string
}

func main() {
	var env Environment
	err := envconfig.Process("", &env)
	if err != nil {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Failed to read enviroment: %q", err)))
		os.Exit(1)
	}

	month := env.Month
	if month == 0 {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Invalid arguments, missing environment variable: 'MONTH'.")))
		os.Exit(1)
	}

	year := env.Year
	if year == 0 {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Invalid arguments, missing environment variable: 'YEAR'.")))
		os.Exit(1)
	}

	outputPath := env.OutputFolder
	if outputPath == "" {
		outputPath = "/output"
	}

	crawlerVersion := env.GitCommit
	if crawlerVersion == "" {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Invalid arguments, missing environment variable: 'GIT_COMMIT'.\n")))
		os.Exit(1)
	}

	if month < 1 || month > 12 {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("Invalid month %d: InvalidParameters.\n", month)))
		os.Exit(1)
	}

	currentYear, currentMonth, _ := time.Now().Date()

	if (year == int(currentYear)) && (year > int(currentMonth)) {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("As master Yoda would say: 'one must not crawl/parse the future %s/%s'.\n", month, year)))
		os.Exit(1)
	}
	if year > int(currentYear) {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("As master Yoda would say: 'one must not crawl/parse the future %s/%s'.\n", month, year)))
		os.Exit(1)
	}

	// Main execution
	fileNames := Crawl(month, year, outputPath)
	employees := Parse(fileNames, month, year)

	cr := coletores.ExecutionResult{
		Cr: coletores.CrawlingResult{
			AgencyID:  "mpro",
			Month:     month,
			Year:      year,
			Files:     fileNames,
			Employees: employees,
			Crawler: coletores.Crawler{
				CrawlerID:      "mpro",
				CrawlerVersion: crawlerVersion,
			},
			Timestamp: time.Now(),
		},
	}

	result, err := json.MarshalIndent(cr, "", "  ")
	if err != nil {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("JSON marshiling error: %q", err)))
		os.Exit(1)
	}
	fmt.Println(string(result))
}
