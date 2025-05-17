package web

import "fmt"

// searchProvider fetches products from a provider for a given query, paginating up to maxPages.
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

// Search coordinates product searches across multiple providers.
type Search struct {
	providers []Provider
}

// SearchResult contains the results and errors from searching a single provider.
type SearchResult struct {
	Errors   []error
	Products []Product
	Provider string
}

// NewSearch creates a new Search instance with the given providers.
func NewSearch(providers []Provider) *Search {
	return &Search{
		providers: providers,
	}
}

// Search runs the query across all providers and returns a slice of SearchResult.
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
