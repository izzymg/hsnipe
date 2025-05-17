package web

type ProductStock string

const (
	StockIn  = "In"
	StockOut = "Out"
)

// Product data
type Product struct {
	Code  string
	Price float64
	Title string
	Stock ProductStock
}
