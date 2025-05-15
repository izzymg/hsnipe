package web

import (
	"fmt"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

func createClient(timeout time.Duration, baseUrl string) *webClient {
	http := &http.Client{
		Timeout: timeout,
	}
	return &webClient{
		client:  http,
		baseUrl: baseUrl,
	}
}

type webClient struct {
	client  *http.Client
	baseUrl string
}

func (w *webClient) get(path string, query map[string]string, expectStatus int) (*html.Node, error) {
	url := fmt.Sprintf("%s/%s", w.baseUrl, path)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()

	res, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != expectStatus {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	doc, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	return doc, nil
}
