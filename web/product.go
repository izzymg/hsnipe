package web

type ProductStock string

const (
	StockIn  = "In"
	StockOut = "Out"
)

type Product struct {
	Code  string
	Price float64
	Title string
	Stock ProductStock
}
