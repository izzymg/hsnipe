package web

// A provider is a plugable search provider, allowing searching a specific product source.
type Provider interface {
	Name() string
	SearchPage(query string, page int) ([]Product, error)
}
