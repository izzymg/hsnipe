package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/izzymg/hsnipe/web"
)

// ExportCSV writes the search results to a CSV file.
func ExportCSV(path string, results []web.SearchResult) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	writer.Write([]string{"Provider", "Title", "Code", "Price", "Error"})

	for _, result := range results {
		errStr := ""
		if len(result.Errors) > 0 {
			var sb strings.Builder
			for _, e := range result.Errors {
				sb.WriteString(e.Error())
				sb.WriteString("; ")
			}
			errStr = sb.String()
		}
		for _, product := range result.Products {
			writer.Write([]string{
				result.Provider,
				product.Title,
				product.Code,
				fmt.Sprintf("%.2f", product.Price),
				errStr,
			})
		}
		// If there are errors but no products, still write the error
		if len(result.Products) == 0 && errStr != "" {
			writer.Write([]string{result.Provider, "", "", "", errStr})
		}
	}
	return writer.Error()
}

// ExportJSON writes the search results to a JSON file.
func ExportJSON(path string, results []web.SearchResult) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Convert errors to string for JSON output
	type ProductOut struct {
		Title string  `json:"title"`
		Code  string  `json:"code"`
		Price float64 `json:"price"`
	}
	type SearchResultOut struct {
		Provider string       `json:"provider"`
		Products []ProductOut `json:"products"`
		Errors   []string     `json:"errors"`
	}
	out := make([]SearchResultOut, 0, len(results))
	for _, r := range results {
		products := make([]ProductOut, 0, len(r.Products))
		for _, p := range r.Products {
			products = append(products, ProductOut{
				Title: p.Title,
				Code:  p.Code,
				Price: p.Price,
			})
		}
		errors := make([]string, 0, len(r.Errors))
		for _, e := range r.Errors {
			errors = append(errors, e.Error())
		}
		out = append(out, SearchResultOut{
			Provider: r.Provider,
			Products: products,
			Errors:   errors,
		})
	}
	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
