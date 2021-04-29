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
	Month        int    `envconfig:"MONTH" required:"true"`
	Year         int    `envconfig:"YEAR" required:"true"`
	OutputFolder string `envconfig:"OUTPUT_FOLDER" default:"/output"`
	GitCommit    string `envconfig:"GIT_COMMIT" required:"true"`
}

func main() {
	var env Environment
	err := envconfig.Process("", &env)
	if err != nil {
		status.ExitFromError(status.NewError(status.DataUnavailable, fmt.Errorf("Failed to read enviroment: %q", err)))
		os.Exit(1)
	}

	month := env.Month
	year := env.Year
	outputPath := env.OutputFolder
	crawlerVersion := env.GitCommit

	if month < 1 || month > 12 {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("Invalid month %d: InvalidParameters.\n", month)))
		os.Exit(1)
	}

	now := time.Now()
	currData := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	crawlDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, now.Location())

	if crawlDate.After(currData) {
		status.ExitFromError(status.NewError(status.SystemError, fmt.Errorf("As master Yoda would say: 'one must not crawl/parse the future %s/%s'.\n", month, year)))
		os.Exit(1)
	}

	// Main execution
	fileNames, err := Crawl(month, year, outputPath)
	if err != nil {
		os.Exit(1)
	}
	employees := Parse(month, year, fileNames)

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
	}
	fmt.Println(string(result))
}
