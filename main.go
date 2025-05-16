package main

import (
	"flag"
	"fmt"
	"os"
	"slices"

	"github.com/izzymg/hsnipe/config"
	"github.com/izzymg/hsnipe/web"
)

func main() {

	configFilePath := "config.json"
	var configFlag = flag.String("config", configFilePath, "Path to the config file")
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
		//web.NewPBTechProvider(),
		web.NewComputerLoungeProvider(),
	})

	fmt.Printf("Searching for %s...\n", config.SearchTerm)

	results, err := search.Search(config.SearchTerm)
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Printf("Results: Provider: %s\n", result.Provider)
		// Sort products on price
		slices.SortFunc(result.Products, func(a, b web.Product) int {
			return int(a.Price - b.Price)
		})

		// Print errors
		if len(result.Errors) > 0 {
			fmt.Printf("Errors: \n")
			for _, err := range result.Errors {
				fmt.Printf("\t%s\n", err)
			}
		}

		// Print products
		for _, product := range result.Products {
			fmt.Printf("%s \t\t%s $%f\n", product.Title, product.Code, product.Price)
		}
	}

}
