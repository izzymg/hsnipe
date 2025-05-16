package parse

import (
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

// GetNodeAttr finds the first attribute with the specified key
// Returns the value of the attribute if found, otherwise returns an empty string.
func GetNodeAttr(n *html.Node, key string) string {
	for _, attr := range n.Attr {
		if attr.Key == key {
			return attr.Val
		}
	}
	return ""
}

// / FindElementByText finds the first text node that contains the specified text
// in the HTML node tree. It performs a depth-first search.
func FindElementByText(n *html.Node, text string) *html.Node {
	for n := range n.Descendants() {
		if n.Type == html.TextNode && strings.Contains(n.Data, text) {
			return n
		}
	}
	return nil
}

// FindElementByClass finds the first element with the specified class name
// in the HTML node tree. It performs a depth-first search.
// It returns nil if no such element is found.
func FindElementByClass(parent *html.Node, class string) *html.Node {
	for n := range parent.Descendants() {
		if n.Type == html.ElementNode {
			if strings.Contains(GetNodeAttr(n, "class"), class) {
				return n
			}
		}
	}
	return nil
}

// FindElementByID finds the first element with the specified ID
// in the HTML node tree. It performs a depth-first search.
func FindElementByID(parent *html.Node, id string) *html.Node {
	for n := range parent.Descendants() {
		if GetNodeAttr(n, "id") == id {
			return n
		}
	}
	return nil
}

func ParsePriceString(text string, prefix bool) (float64, error) {
	if prefix {
		cut, ok := strings.CutPrefix(text, "$")
		if !ok {
			return -1, fmt.Errorf("unexpected full price format")
		}
		text = cut
	}
	text = strings.ReplaceAll(text, ",", "")
	text = strings.ReplaceAll(text, "\n", "")
	formattedPrice, err := strconv.ParseFloat(text, 32)
	if err != nil {
		return -1, fmt.Errorf("failed to parse full price as float")
	}
	return formattedPrice, nil
}
