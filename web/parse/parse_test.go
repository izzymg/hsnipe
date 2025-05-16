package parse_test

import (
	"strings"
	"testing"

	"github.com/izzymg/hsnipe/web/parse"
	"golang.org/x/net/html"
)

func getSampleHTML() string {
	return `
	<html>
		<head>
			<title>Test</title>
		</head>
		<body>
			<div class="product">
				<a href="https://example.com/product1" class="product-link">Product 1</a>
			</div>
			<div class="product">
				<a href="https://example.com/product2" id="product-2" class="product-link">Product 2</a>
			</div>
			<div class="product">
				<a href="https://example.com/product3" class="product-link">Product 3</a>
			</div>
		</body>
	</html>`
}

func TestFindAttr(t *testing.T) {
	htmlContent := getSampleHTML()
	node, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find the first product link
	productLinkNode := parse.FindElementByClass(node, "product-link")
	if productLinkNode == nil {
		t.Fatal("Failed to find product link node")
	}

	// Get the href attribute
	href := parse.GetNodeAttr(productLinkNode, "href")
	expectedHref := "https://example.com/product1"
	if href != expectedHref {
		t.Errorf("Expected href %s, got %s", expectedHref, href)
	}
}

func TestFindElementByClass(t *testing.T) {
	htmlContent := getSampleHTML()
	node, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find the first product node
	productNode := parse.FindElementByClass(node, "product")
	if productNode == nil {
		t.Fatal("Failed to find product node")
	}
}

func TestFindByText(t *testing.T) {
	htmlContent := getSampleHTML()
	node, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// Find the first product link by text
	productLinkNode := parse.FindElementByText(node, "Product 1")
	if productLinkNode == nil {
		t.Fatal("Failed to find product link node by text")
	}

	// Get the href attribute
	href := parse.GetNodeAttr(productLinkNode.Parent, "href")
	expectedHref := "https://example.com/product1"
	if href != expectedHref {
		t.Errorf("Expected href %s, got %s", expectedHref, href)
	}
}

func TestFindElementByID(t *testing.T) {
	htmlContent := getSampleHTML()
	node, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		t.Fatalf("Failed to parse HTML: %v", err)
	}

	// No ID
	productNode := parse.FindElementByID(node, "product-1")
	if productNode != nil {
		t.Fatal("Expected to not find product node by ID, but found one")
	}
	// Some ID
	productNode = parse.FindElementByID(node, "product-2")
	if productNode == nil {
		t.Fatal("Expected to find product node by ID, but did not")
	}

}

func TestParsePriceString(t *testing.T) {
	out, err := parse.ParsePriceString("$1,899.00", true)
	if err != nil {
		t.Fatalf("got error parsing price %v", err)
	}
	if out != 1899.00 {
		t.Fatalf("expected %f, got %f", 1899.00, out)
	}

	out, err = parse.ParsePriceString("1,899.00", false)
	if err != nil {
		t.Fatalf("got error parsing price %v", err)
	}
	if out != 1899.00 {
		t.Fatalf("expected %f, got %f", 1899.00, out)
	}
}
