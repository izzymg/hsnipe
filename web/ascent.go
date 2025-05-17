package web

import (
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/izzymg/hsnipe/web/client"
)

type ascentSearchRequest struct {
	QueryBy             string `json:"query_by"`
	SplitJoinTokens     string `json:"split_join_tokens"`
	HighlightFullFields string `json:"highlight_full_fields"`
	Collection          string `json:"collection"`
	Q                   string `json:"q"`
	FacetBy             string `json:"facet_by"`
	MaxFacetValues      int    `json:"max_facet_values"`
	Page                int    `json:"page"`
	PerPage             int    `json:"per_page"`
}

type ascentSearchResponse struct {
	Results []ascentSearchResult `json:"results"`
}

type ascentSearchResult struct {
	FacetCounts   []interface{}          `json:"facet_counts"`
	Found         int                    `json:"found"`
	Hits          []ascentSearchHit      `json:"hits"`
	OutOf         int                    `json:"out_of"`
	Page          int                    `json:"page"`
	RequestParams map[string]interface{} `json:"request_params"`
	SearchCutoff  bool                   `json:"search_cutoff"`
	SearchTimeMs  int                    `json:"search_time_ms"`
}

type ascentSearchHit struct {
	Document   ascentProductDocument      `json:"document"`
	Highlight  map[string]ascentHighlight `json:"highlight"`
	Highlights []ascentHighlight          `json:"highlights"`
	// text_match and text_match_info omitted for brevity
}

type ascentProductDocument struct {
	AscentPartNumber       int                   `json:"AscentPartNumber"`
	AscentPartNumberSearch string                `json:"AscentPartNumberSearch"`
	Brand                  string                `json:"Brand"`
	Category               string                `json:"Category"`
	Chipsetmanufacturer    string                `json:"Chipsetmanufacturer"`
	DVIDports              int                   `json:"DVIDports"`
	DateCreated            int64                 `json:"DateCreated"`
	IsBusinessProduct      bool                  `json:"IsBusinessProduct"`
	LastUpdated            int64                 `json:"LastUpdated"`
	ManufacturerPartNumber string                `json:"ManufacturerPartNumber"`
	Memory                 int                   `json:"Memory"`
	MiniDisplayPort        int                   `json:"MiniDisplayPort"`
	Popularity             int                   `json:"Popularity"`
	Price                  float64               `json:"Price"`
	PriceInclGST           float64               `json:"PriceInclGST"`
	ProductName            string                `json:"ProductName"`
	Productline            string                `json:"Productline"`
	RAMtype                string                `json:"RAMtype"`
	RRP                    int                   `json:"RRP"`
	Stock                  string                `json:"Stock"`
	SupplierPartNumber     string                `json:"SupplierPartNumber"`
	Tag                    []interface{}         `json:"Tag"`
	TotalDisplayPortports  int                   `json:"TotalDisplayPortports"`
	TotalHDMIports         int                   `json:"TotalHDMIports"`
	TotalVGAports          int                   `json:"TotalVGAports"`
	Type                   []string              `json:"Type"`
	Videocardchipset       string                `json:"Videocardchipset"`
	Videocardchipsetseries string                `json:"Videocardchipsetseries"`
	ID                     string                `json:"id"`
	ImageLinks             []ascentImageLink     `json:"imageLinks"`
	MarketingSpec          []ascentMarketingSpec `json:"marketingSpec"`
	MediaLink              []interface{}         `json:"mediaLink"`
	TagCollection          []interface{}         `json:"tagCollection"`
}

type ascentImageLink struct {
	ImageOrder int    `json:"imageOrder"`
	ImageURL   string `json:"imageURL"`
}

type ascentMarketingSpec struct {
	Priority bool   `json:"priority"`
	Source   string `json:"source"`
	Spec     string `json:"spec"`
}

type ascentHighlight struct {
	Field         string   `json:"field,omitempty"`
	MatchedTokens []string `json:"matched_tokens"`
	Snippet       string   `json:"snippet"`
	Value         string   `json:"value"`
}

type AscentProvider struct {
	webClient *client.WebClient
	apiClient *client.WebClient
}

func NewAscentProvider() *AscentProvider {
	return &AscentProvider{
		webClient: client.CreateClient(time.Second*10, "https://ascent.co.nz"),
		apiClient: client.CreateClient(time.Second*10, "https://83rynw1spubc6vg0p-1.a1.typesense.net"),
	}
}

func (ap AscentProvider) Name() string {
	return "Ascent"
}

func (ap AscentProvider) SearchPage(query string, page int) ([]Product, error) {
	js, err := ap.webClient.GetRaw("search/app5.js", map[string]string{}, 200)
	if err != nil {
		return nil, err
	}

	regex := regexp.MustCompile(`apiKey: "(.*)",`)
	match := regex.FindSubmatch(js)
	if match == nil || len(match) != 2 {
		return nil, fmt.Errorf("failed to find key")
	}
	key := string(match[1])

	response, err := ap.apiClient.PostJson("multi_search", map[string]string{
		"x-typesense-api-key": key,
	}, map[string][]ascentSearchRequest{
		"searches": {
			ascentSearchRequest{
				QueryBy:         "AscentPartNumberSearch,Brand,ProductName,ManufacturerPartNumber,SupplierPartNumber,Category",
				Collection:      "Products",
				FacetBy:         "Category",
				Page:            page,
				PerPage:         250,
				Q:               query,
				MaxFacetValues:  100,
				SplitJoinTokens: "always",
			},
		},
	}, 200)
	if err != nil {
		return nil, err
	}

	responseObj := &ascentSearchResponse{}
	json.Unmarshal(response, responseObj)

	var products []Product
	if len(responseObj.Results) > 0 {
		for _, hit := range responseObj.Results[0].Hits {
			doc := hit.Document
			stock := StockOut
			if doc.Stock == "In stock" {
				stock = StockIn
			} else {
				stock = StockOut
			}
			products = append(products, Product{
				Code:  doc.ID,
				Price: doc.PriceInclGST,
				Title: doc.ProductName,
				Stock: ProductStock(stock),
			})
		}
	}

	return products, nil
}
