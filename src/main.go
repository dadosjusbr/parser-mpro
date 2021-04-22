package main

import (
    "fmt";
    "os";
    "time";
    "strconv";
    "encoding/json";
    "github.com/dadosjusbr/coletores"
)

func main(){
    month, ok := os.LookupEnv("MONTH")
    if !ok {
        fmt.Fprintln(os.Stderr, "Invalid arguments, missing environment variable: 'MONTH'.\n")
        os.Exit(1)
    }
    
    year, ok := os.LookupEnv("YEAR")
    if !ok {
        fmt.Fprintln(os.Stderr, "Invalid arguments, missing environment variable: 'YEAR'.\n")
        os.Exit(1)
    }
    
    output_path, ok := os.LookupEnv("OUTPUT_FOLDER")
    if !ok{
        output_path = "/output"
    }
    
    crawler_version, ok :=  os.LookupEnv("GIT_COMMIT")
    if !ok{
        fmt.Fprintln(os.Stderr, "Invalid arguments, missing environment variable: 'GIT_COMMIT'.\n")
        os.Exit(1)
    }
    
    current_year, current_month, _ := time.Now().Date()
    
    month_check, err := strconv.Atoi(month)
    if err == nil{
        if ( month_check < 1 || month_check > 12){
            fmt.Fprintln(os.Stderr, "Invalid month %s: InvalidParameters.\n", month)
            os.Exit(1)
        }
    } 

    year_check, err := strconv.Atoi(year)
    if err == nil{
        if ((year_check == int(current_year)) && ( year_check > int(current_month))){
            fmt.Fprintln(os.Stderr, "As master Yoda would say: 'one must not crawl/parse the future %s/%s'.\n", month, year)
            os.Exit(1)
        }
        if (year_check > int(current_year)){
            fmt.Fprintln(os.Stderr, "As master Yoda would say: 'one must not crawl/parser th future %s/%s'.\n", month, year)
        }
    }

    // Main execution
    file_names := Crawl(month, year, output_path)
    employees := Parse(file_names, month, year)

    cr := coletores.ExecutionResult{
        Cr: coletores.CrawlingResult{
            AgencyID:  "mppe",
            Month:     month_check,
            Year:      year_check,
            Files:     file_names,
            Employees: employees,
            Crawler: coletores.Crawler{
                CrawlerID:      "mppe",
                CrawlerVersion: crawler_version,
            },
            Timestamp: time.Now(),
        },
    }

    result, err := json.MarshalIndent(cr, "", "  ")
	if err != nil {
		fmt.Fprintln(os.Stderr,"JSON marshiling error: %q", err)
	}
	fmt.Println(string(result))
}




