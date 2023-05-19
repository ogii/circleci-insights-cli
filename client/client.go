package client

import (
	"net/http"
	"time"
)

type Client struct {
	HTTPClient *http.Client
	BaseURL    string
	Token      string
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
		BaseURL: baseURL,
		Token:   token,
	}
}
