package client

import (
	"net/http"
)

type HTTPClient struct {
	baseURL string
	Client  *http.Client
}

func NewHTTPClient(baseURL string, client *http.Client) *HTTPClient {
	return &HTTPClient{baseURL: baseURL, Client: client}
}
