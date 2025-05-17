package main

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/izzymg/hsnipe/config"
	"github.com/izzymg/hsnipe/web"
)

func main() {
	configFilePath := "config.json"
	var configFlag = flag.String("config", configFilePath, "Path to the config file")
	var exportPath = flag.String("export", "", "Path to export results (CSV or JSON)")
	var exportFormat = flag.String("format", "csv", "Export format: csv or json")
	flag.Parse()
	if configFlag != nil {
		configFilePath = *configFlag
	}

	config, err := config.ParseConfig(configFilePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	search := web.NewSearch([]web.Provider{
		web.NewPBTechProvider(*regexp.MustCompile(config.PBTechConfig.Filter)),
		web.NewComputerLoungeProvider(*regexp.MustCompile(config.ComputerLoungeConfig.TitleFilter)),
		web.NewAscentProvider(),
	})

	fmt.Printf("Searching for %s...\n", config.SearchTerm)

	results, err := search.Search(config.SearchTerm)
	if err != nil {
		panic(err)
	}

	// Sort products by price for each provider
	for i := range results {
		slices.SortFunc(results[i].Products, func(a, b web.Product) int {
			return int(a.Price - b.Price)
		})
	}

	// Export if requested
	if *exportPath != "" {
		switch strings.ToLower(*exportFormat) {
		case "csv":
			if err := ExportCSV(*exportPath, results); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export CSV: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Exported results to %s (CSV)\n", *exportPath)
		case "json":
			if err := ExportJSON(*exportPath, results); err != nil {
				fmt.Fprintf(os.Stderr, "Failed to export JSON: %v\n", err)
				os.Exit(1)
			}
			fmt.Printf("Exported results to %s (JSON)\n", *exportPath)
		default:
			fmt.Fprintf(os.Stderr, "Unknown export format: %s\n", *exportFormat)
			os.Exit(1)
		}
		return
	}

	// Print to stdout as before
	for _, result := range results {
		fmt.Printf("Results: Provider: %s\n", result.Provider)
		if len(result.Errors) > 0 {
			fmt.Printf("Errors: \n")
			for _, err := range result.Errors {
				fmt.Printf("\t%s\n", err)
			}
		}
		for _, product := range result.Products {
			fmt.Printf("%s \t\t%s $%f\n", product.Title, product.Code, product.Price)
		}
	}
}
