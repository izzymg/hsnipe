package web

type Search struct {
	providers []Provider
}

type SearchResult struct {
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
		products, err := provider.Search(query)
		if err != nil {
			return nil, err
		}

		results = append(results, SearchResult{
			Products: products,
			Provider: provider.Name(),
		})
	}

	return results, nil
}
