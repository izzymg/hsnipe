package main

import (
	"fmt"
	"slices"

	"github.com/izzymg/hsnipe/web"
)

func main() {
	products, err := web.Search("RTX 5070 ti")
	if err != nil {
		panic(err)
	}

	// Sort products on price
	slices.SortFunc(products, func(a, b web.Product) int {
		return int(a.Price - b.Price)
	})

	// Print products
	for _, product := range products {
		fmt.Printf("%s \t\t$%s %f\n", product.Title, product.Code, product.Price)
	}

}
