package web

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/net/html"
)

/** Web: pbtech.co.nz */
const DOMAIN = "https://www.pbtech.co.nz"
const HTTP_TIMEOUT = 10 * time.Second

var PRODUCT_CODE_FILTER = regexp.MustCompile(`^VGA.*$`)

type Product struct {
	Code  string
	Price float64
	Title string
}

func findAttr(n *html.Node, key string) string {
	if n.Type == html.ElementNode {
		for _, attr := range n.Attr {
			if attr.Key == key {
				return attr.Val
			}
		}
	}
	return ""
}

func findByText(n *html.Node, text string) *html.Node {
	for n := range n.Descendants() {
		if n.Type == html.TextNode && strings.Contains(n.Data, text) {
			return n
		}
	}
	return nil
}

// findElementByClass finds the first element with the specified class name
// in the HTML node tree. It performs a depth-first search.
// It returns nil if no such element is found.
func findElementByClass(parent *html.Node, class string) *html.Node {
	for n := range parent.Descendants() {
		if n.Type == html.ElementNode {
			if strings.Contains(findAttr(n, "class"), class) {
				return n
			}
		}
	}
	return nil
}

func parseProductCard(parent *html.Node) (*Product, error) {

	priceContainer := findElementByClass(parent, "item-price-amount")
	if priceContainer == nil {
		return nil, fmt.Errorf("no price container found")
	}
	ginc := findElementByClass(priceContainer, "ginc")
	if ginc == nil {
		return nil, fmt.Errorf("no GST inclusive container found")
	}

	fullPriceEle := findElementByClass(ginc, "full-price")
	if fullPriceEle == nil {
		return nil, fmt.Errorf("no full price element found")
	}

	linkEle := findElementByClass(parent, "js-product-link")
	if linkEle == nil {
		return nil, fmt.Errorf("no link element found")
	}

	titleEle := findElementByClass(parent, "np_title")
	if titleEle == nil {
		return nil, fmt.Errorf("no title element found")
	}

	productCode := findAttr(linkEle, "data-product-code")
	fullPrice, ok := strings.CutPrefix(fullPriceEle.FirstChild.Data, "$")
	if !ok {
		return nil, fmt.Errorf("unexpected full price format")
	}
	fullPrice = strings.ReplaceAll(fullPrice, ",", "")
	formattedPrice, err := strconv.ParseFloat(fullPrice, 32)
	if err != nil {
		return nil, fmt.Errorf("failed to parse full price as float")
	}
	return &Product{
		Code:  productCode,
		Price: formattedPrice,
		Title: strings.Trim(titleEle.FirstChild.Data, " \n"),
	}, nil
}

func parsePage(term string, page int) ([]Product, error) {
	url := fmt.Sprintf("%s/search?sf=%s&pg=%d", DOMAIN, url.QueryEscape(term), page)
	client := &http.Client{
		Timeout: HTTP_TIMEOUT,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	if findByText(doc, "No products") != nil {
		return nil, nil
	}

	n := findElementByClass(doc, "products-view")
	if n == nil {
		return nil, fmt.Errorf("no product view found for %s", term)
	}
	wrapper := n.FirstChild.NextSibling

	products := make([]Product, 0)

	for card := wrapper.FirstChild; card != nil; card = card.NextSibling {
		if card.Type == html.ElementNode && card.Data == "div" {
			if !strings.Contains(findAttr(card, "class"), "js-product-card") {
				continue
			}
			product, err := parseProductCard(card)
			if err != nil {
				return nil, err
			}

			if !PRODUCT_CODE_FILTER.MatchString(product.Code) {
				continue
			}
			products = append(products, *product)
		}
	}
	return products, nil
}

func Search(term string) ([]Product, error) {
	pageLimit := 10
	products := make([]Product, 0)
	for n := 1; n <= pageLimit; n++ {
		fmt.Printf("Searching page %d...\n", n)
		foundProducts, err := parsePage(term, n)
		if err != nil {
			return nil, err
		}
		if len(foundProducts) == 0 {
			fmt.Println("No products were found")
			break
		}

		products = append(products, foundProducts...)
	}
	fmt.Println("Search completed.")
	return products, nil
}
