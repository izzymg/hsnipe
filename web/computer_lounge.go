package web

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/izzymg/hsnipe/web/client"
	"github.com/izzymg/hsnipe/web/parse"
	"golang.org/x/net/html"
)

/** Web: computerlounge.co.nz */

const category = "293914804339"

func NewComputerLoungeProvider(titleFilter regexp.Regexp) *ComputerLoungeProvider {
	return &ComputerLoungeProvider{
		client:      client.CreateClient(10*time.Second, "https://www.computerlounge.co.nz"),
		titleFilter: titleFilter,
	}
}

type ComputerLoungeProvider struct {
	client      *client.WebClient
	titleFilter regexp.Regexp
}

func (c ComputerLoungeProvider) Name() string {
	return "Computer Lounge"
}

func (c ComputerLoungeProvider) parseCard(node *html.Node) (*Product, error) {
	priceEle := parse.FindElementByClass(node, "price__current")
	if priceEle == nil {
		return nil, fmt.Errorf("failed to find price container")
	}
	formattedPrice, err := parse.ParsePriceString(priceEle.FirstChild.Data, true)
	if err != nil {
		return nil, err
	}

	nameEle := parse.FindElementByClass(node, "js-prod-link")
	if nameEle == nil {
		return nil, fmt.Errorf("failed to find product link")
	}
	name := parse.GetNodeAttr(nameEle, "aria-label")

	skuTitleEle := parse.FindElementByText(node, "SKU")
	// SKU sometimes not set.
	code := "Unknown"
	if skuTitleEle != nil {
		code = skuTitleEle.Parent.Parent.NextSibling.NextSibling.FirstChild.Data
	}

	return &Product{
		Price: formattedPrice,
		Code:  code,
		Title: name,
	}, nil
}

func (c ComputerLoungeProvider) SearchPage(query string, page int) ([]Product, error) {
	node, err := c.client.GetHtml("search", map[string]string{
		"q":    query,
		"page": strconv.Itoa(page),
	}, 200)

	if err != nil {
		return nil, err
	}

	wrapperElement := parse.FindElementByClass(node, "main-products-grid__results")

	if wrapperElement == nil {
		return nil, fmt.Errorf("No product container")
	}

	wrapperElement = wrapperElement.FirstChild

	products := make([]Product, 0)
	for cardElement := wrapperElement.FirstChild; cardElement != nil; cardElement = cardElement.NextSibling {
		if !strings.Contains(parse.GetNodeAttr(cardElement, "class"), "js-pagination-result") {
			continue
		}
		productCardElement := cardElement.FirstChild
		if productCardElement == nil || !strings.Contains(parse.GetNodeAttr(productCardElement, "class"), "card--product") {
			continue
		}
		product, err := c.parseCard(productCardElement)
		if err != nil {
			return nil, err
		}
		if !c.titleFilter.MatchString(product.Title) {
			continue
		}
		products = append(products, *product)
	}
	return products, nil
}
