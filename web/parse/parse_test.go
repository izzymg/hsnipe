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
				<a href="https://example.com/product2" class="product-link">Product 2</a>
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
