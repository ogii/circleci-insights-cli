package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRequest(t *testing.T) {
	client := NewClient("http://test-url.com", "000000000000000000000000000000")
	req, err := client.BuildRequest("test-slug", "test-url", "test-branch", "test-window", "")

	assert.NoError(t, err)
	assert.NotNil(t, req)
	assert.Equal(t, "http://test-url.com/insights/test-slug/workflows/test-url?branch=test-branch&reporting-window=test-window", req.URL.String())
	assert.Equal(t, "000000000000000000000000000000", req.Header.Get("Circle-Token"))
}

func TestExecuteRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{ "items" : [] }`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "000000000000000000000000000000")
	req, _ := client.BuildRequest("test-slug", "test-url", "test-branch", "test-window", "")
	resp, err := client.ExecuteRequest(req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHandleResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{ "next_page_token" : null, "items" : [] }`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "000000000000000000000000000000")
	req, _ := client.BuildRequest("test-slug", "test-url", "test-branch", "test-window", "")
	resp, _ := client.ExecuteRequest(req)
	insightsSummary, err := client.HandleResponse(resp)

	assert.NoError(t, err)
	assert.NotNil(t, insightsSummary)
	assert.Equal(t, 0, len(insightsSummary.Workflows))
}

func TestFetchInsightsSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{ "next_page_token" : null, "items" : [] }`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "000000000000000000000000000000")
	insightsSummary, err := client.FetchInsightsSummary("test-slug", "test-url", "test-branch", "test-window")

	assert.NoError(t, err)
	assert.NotNil(t, insightsSummary)
	assert.Equal(t, 0, len(insightsSummary.Workflows))
}
