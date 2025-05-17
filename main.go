package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/izzymg/hsnipe/config"
	"github.com/izzymg/hsnipe/web"
)

func main() {
	var quietFlag = flag.Bool("quiet", false, "No stdout")
	configFilePath := "config.json"
	var configFlag = flag.String("config", configFilePath, "Path to the config file")
	var exportPath = flag.String("export", "", "Path to export results (CSV or JSON)")
	var exportFormat = flag.String("format", "csv", "Export format: csv or json")
	flag.Parse()
	if configFlag != nil {
		configFilePath = *configFlag
	}

	logger := log.New(os.Stdout, "hsnipe: ", log.Lmsgprefix)
	if *quietFlag {
		logger.SetOutput(io.Discard)
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
	}, logger)

	logger.Printf("Searching for %s...\n", config.SearchTerm)

	results, err := search.Search(config.SearchTerm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
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
				logger.Fatalf("Failed to export CSV: %v\n", err)
				os.Exit(1)
			}
			logger.Printf("Exported results to %s (CSV)\n", *exportPath)
		case "json":
			if err := ExportJSON(*exportPath, results); err != nil {
				logger.Fatalf("Failed to export JSON: %v\n", err)
				os.Exit(1)
			}
			logger.Printf("Exported results to %s (JSON)\n", *exportPath)
		default:
			logger.Fatalf("Unknown export format: %s\n", *exportFormat)
			os.Exit(1)
		}
		return
	}

	// Print to stdout as before
	for _, result := range results {
		logger.Printf("Results: Provider: %s\n", result.Provider)
		if len(result.Errors) > 0 {
			logger.Printf("Errors: \n")
			for _, err := range result.Errors {
				logger.Printf("\t%s\n", err)
			}
		}
		for _, product := range result.Products {
			logger.Printf("%s \t\t%s $%f\n", product.Title, product.Code, product.Price)
		}
	}
}
