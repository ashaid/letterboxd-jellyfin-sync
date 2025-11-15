package utils

import (
	"fmt"
	"io"
	"net/http"
)

type SimpleClient struct {
	BaseURL string
	Client  *http.Client
}

func NewSimpleClient(baseURL string) *SimpleClient {
	return &SimpleClient{
		BaseURL: baseURL,
		Client:  &http.Client{},
	}
}

func (dc *SimpleClient) Get(path string) ([]byte, error) {
	url := dc.BaseURL + path

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	SetWebHeaders(req)

	res, err := dc.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("request failed with status %d", res.StatusCode)
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return body, nil
}
