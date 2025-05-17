package web

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/html"
)

const debug = false

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

func (w *webClient) buildRequest(path string, query map[string]string, method string, body io.Reader) (*http.Request, error) {
	url := fmt.Sprintf("%s/%s", w.baseUrl, path)
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	for key, value := range query {
		q.Add(key, value)
	}
	req.URL.RawQuery = q.Encode()
	return req, nil
}

func (w *webClient) postJson(path string, query map[string]string, body any, expectStatus int) ([]byte, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	jsonReq := bytes.NewBuffer(data)
	req, err := w.buildRequest(path, query, "POST", jsonReq)
	res, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != expectStatus {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	obj, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (w *webClient) getRaw(path string, query map[string]string, expectStatus int) ([]byte, error) {
	req, err := w.buildRequest(path, query, "GET", nil)
	if err != nil {
		return nil, err
	}
	res, err := w.client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != expectStatus {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (w *webClient) getHtml(path string, query map[string]string, expectStatus int) (*html.Node, error) {
	req, err := w.buildRequest(path, query, "GET", nil)
	if err != nil {
		return nil, err
	}
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
