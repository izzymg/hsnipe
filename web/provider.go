package web

type Provider interface {
	Name() string
	Search(query string) ([]Product, error)
}
