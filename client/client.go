package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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

func (c *Client) BuildRequest(slug, urlPath, branch, reportingWindow, nextPageToken string) (*http.Request, error) {
	baseURL, err := url.Parse(c.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %v", err)
	}

	// Construct the path
	path := fmt.Sprintf("/api/v2/insights/%s/workflows/%s", slug, urlPath)

	// Construct the query parameters
	params := url.Values{}
	params.Add("branch", branch)
	params.Add("reporting-window", reportingWindow)
	if nextPageToken != "" {
		params.Add("next_page_token", nextPageToken)
	}

	// Construct the final URL
	finalURL := baseURL.ResolveReference(&url.URL{Path: path, RawQuery: params.Encode()})

	req, err := http.NewRequest("GET", finalURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request for URL %s: %v", finalURL, err)
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
