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
		req, err := c.BuildRequest(slug, url, branch, reportingWindow, nextPageToken)
		if err != nil {
			return nil, err
		}

		resp, err := c.ExecuteRequest(req)
		if err != nil {
			return nil, err
		}

		currentPage, err := c.HandleResponse(resp)
		if err != nil {
			return nil, err
		}

		insightsSummary.Workflows = append(insightsSummary.Workflows, currentPage.Workflows...)

		if currentPage.NextPageToken == "" {
			break
		}

		nextPageToken = currentPage.NextPageToken
	}

	return &insightsSummary, nil
}

func (c *Client) BuildRequest(slug, url, branch, reportingWindow, nextPageToken string) (*http.Request, error) {
	url = fmt.Sprintf("%s/insights/%s/workflows/%s?branch=%s&reporting-window=%s", c.BaseURL, slug, url, branch, reportingWindow)
	if nextPageToken != "" {
		url += "&next_page_token=" + nextPageToken
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request for URL %s: %v", url, err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", c.Token)

	return req, nil
}

func (c *Client) ExecuteRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request to URL %s: %v", req.URL, err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("non-200 response from CircleCI API for URL %s: %s", req.URL, resp.Status)
	}
	return resp, nil
}

func (c *Client) HandleResponse(resp *http.Response) (*data.InsightsSummary, error) {
	defer resp.Body.Close()

	var insightsSummary data.InsightsSummary
	err := json.NewDecoder(resp.Body).Decode(&insightsSummary)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response from URL %s: %v", resp.Request.URL, err)
	}
	return &insightsSummary, nil
}
