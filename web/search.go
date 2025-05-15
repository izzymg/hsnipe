package web

import "fmt"

func searchProvider(provider Provider, query string) ([]Product, error) {
	maxPages := 10
	products := make([]Product, 0)
	for page := 1; page <= maxPages; page++ {
		fmt.Printf("%s: Searching page %d...\n", provider.Name(), page)
		foundProducts, err := provider.SearchPage(query, page)
		if err != nil {
			return nil, err
		}
		if len(foundProducts) == 0 {
			fmt.Printf("%s: Reached page %d, no results.\n", provider.Name(), page)
			break
		}
		products = append(products, foundProducts...)
	}
	return products, nil
}

type Search struct {
	providers []Provider
}

type SearchResult struct {
	Errors   []error
	Products []Product
	Provider string
}

func NewSearch(providers []Provider) *Search {
	return &Search{
		providers: providers,
	}
}

func (search Search) Search(query string) ([]SearchResult, error) {
	results := make([]SearchResult, 0, len(search.providers))

	for _, provider := range search.providers {
		errors := make([]error, 0, len(search.providers))
		products, err := searchProvider(provider, query)
		if err != nil {
			errors = append(errors, fmt.Errorf("error searching %s: %w", provider.Name(), err))
		}
		results = append(results, SearchResult{
			Products: products,
			Provider: provider.Name(),
			Errors:   errors,
		})
	}

	return results, nil
}
