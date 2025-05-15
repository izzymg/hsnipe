package main

import (
	"fmt"
	"slices"

	"github.com/izzymg/hsnipe/web"
)

func main() {

	search := web.NewSearch([]web.Provider{
		web.PBTechProvider{},
	})

	results, err := search.Search("RTX 5070 ti")
	if err != nil {
		panic(err)
	}

	for _, result := range results {
		fmt.Printf("Provider: %s\n", result.Provider)
		// Sort products on price
		slices.SortFunc(result.Products, func(a, b web.Product) int {
			return int(a.Price - b.Price)
		})

		// Print products
		for _, product := range result.Products {
			fmt.Printf("%s \t\t$%s %f\n", product.Title, product.Code, product.Price)
		}
	}

}
