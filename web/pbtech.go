package web

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/izzymg/hsnipe/web/parse"
	"golang.org/x/net/html"
)

/** Web: pbtech.co.nz */

var productFilter = regexp.MustCompile(`^VGA.*$`)

type PBTechProvider struct {
	client *webClient
}

func NewPBTechProvider() *PBTechProvider {
	return &PBTechProvider{
		client: createClient(10*time.Second, "https://www.pbtech.co.nz"),
	}
}

func (p PBTechProvider) parseProductCard(parent *html.Node) (*Product, error) {

	priceContainer := parse.FindElementByClass(parent, "item-price-amount")
	if priceContainer == nil {
		return nil, fmt.Errorf("no price container found")
	}
	ginc := parse.FindElementByClass(priceContainer, "ginc")
	if ginc == nil {
		return nil, fmt.Errorf("no GST inclusive container found")
	}

	fullPriceEle := parse.FindElementByClass(ginc, "full-price")
	if fullPriceEle == nil {
		return nil, fmt.Errorf("no full price element found")
	}

	linkEle := parse.FindElementByClass(parent, "js-product-link")
	if linkEle == nil {
		return nil, fmt.Errorf("no link element found")
	}

	titleEle := parse.FindElementByClass(parent, "np_title")
	if titleEle == nil {
		return nil, fmt.Errorf("no title element found")
	}

	productCode := parse.GetNodeAttr(linkEle, "data-product-code")
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

func (p PBTechProvider) searchPage(term string, page int) ([]Product, error) {
	doc, err := p.client.get("search", map[string]string{
		"sf": term,
		"pg": strconv.Itoa(page),
	}, 200)
	if err != nil {
		return nil, err
	}

	if parse.FindElementByText(doc, "No products") != nil {
		return nil, nil
	}

	n := parse.FindElementByClass(doc, "products-view")
	if n == nil {
		return nil, fmt.Errorf("no product view found for %s", term)
	}
	wrapper := n.FirstChild.NextSibling

	products := make([]Product, 0)

	for card := wrapper.FirstChild; card != nil; card = card.NextSibling {
		if card.Type == html.ElementNode && card.Data == "div" {
			if !strings.Contains(parse.GetNodeAttr(card, "class"), "js-product-card") {
				continue
			}
			product, err := p.parseProductCard(card)
			if err != nil {
				return nil, err
			}

			if !productFilter.MatchString(product.Code) {
				continue
			}
			products = append(products, *product)
		}
	}
	return products, nil
}

func (p PBTechProvider) Name() string {
	return "PBTech"
}

func (p PBTechProvider) SearchPage(term string, page int) ([]Product, error) {
	foundProducts, err := p.searchPage(term, page)
	if err != nil {
		return nil, err
	}
	return foundProducts, nil
}
