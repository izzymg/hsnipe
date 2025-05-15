package web

type Provider interface {
	Name() string
	SearchPage(query string, page int) ([]Product, error)
}
