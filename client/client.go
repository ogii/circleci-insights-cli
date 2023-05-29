package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ogii/circleci-insights-cli/data"
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

func (c *Client) FetchInsightsSummary(slug, url, branch, reportingWindow string) (*data.InsightsSummary, error) {
	var insightsSummary data.InsightsSummary
	var nextPageToken string

	for {
		url := fmt.Sprintf("%s/insights/%s/workflows/%s?branch=%s&reporting-window=%s", c.BaseURL, slug, url, branch, reportingWindow)
		fmt.Println("View in browser: " + url)
		if nextPageToken != "" {
			url += "&page-token=" + nextPageToken
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request for URL %s: %v", url, err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Circle-Token", c.Token)

		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error making HTTP request to URL %s: %v", url, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("non-200 response from CircleCI API for URL %s: %s", url, resp.Status)
		}

		var currentPage data.InsightsSummary
		err = json.NewDecoder(resp.Body).Decode(&currentPage)
		if err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("error unmarshaling JSON response from URL %s: %v", url, err)
		}

		insightsSummary.Workflows = append(insightsSummary.Workflows, currentPage.Workflows...)

		if currentPage.NextPageToken == "" {
			break
		}

		nextPageToken = currentPage.NextPageToken
	}
	return &insightsSummary, nil
}
