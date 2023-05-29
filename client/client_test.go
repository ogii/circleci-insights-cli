package client

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchInsightsSummary(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, `{
			"next_page_token" : null,
			"items" : [ {
			  "name" : "test-workflow",
			  "metrics" : {
				"total_runs" : 54,
				"successful_runs" : 52,
				"mttr" : 239676,
				"total_credits_used" : 34863,
				"failed_runs" : 1,
				"median_credits_used" : 0,
				"success_rate" : 0.962962963,
				"duration_metrics" : {
				  "min" : 69,
				  "mean" : 356,
				  "median" : 318,
				  "p95" : 451,
				  "max" : 1712,
				  "standard_deviation" : 203.0,
				  "total_duration" : 0
				},
				"total_recoveries" : 0,
				"throughput" : 0.6
			  },
			  "window_start" : "2023-02-17T14:43:26.205Z",
			  "window_end" : "2023-05-10T13:56:29.813Z",
			  "project_id" : "7097f60c-74d1-4936-8d1a-268d4042a493"
			} ]
		  }`)
	}))
	defer server.Close()

	client := NewClient(server.URL, "test-token")
	insightsSummary, err := client.FetchInsightsSummary("test-slug", "test-url", "test-branch", "test-window")

	assert.NoError(t, err)
	assert.NotNil(t, insightsSummary)
	assert.Equal(t, 1, len(insightsSummary.Workflows))
	assert.Equal(t, "test-workflow", insightsSummary.Workflows[0].Name)
	assert.Equal(t, 0.962962963, insightsSummary.Workflows[0].Metrics.SuccessRate)
	assert.Equal(t, 54, insightsSummary.Workflows[0].Metrics.TotalRuns)
	assert.Equal(t, 1, insightsSummary.Workflows[0].Metrics.FailedRuns)
	assert.Equal(t, 52, insightsSummary.Workflows[0].Metrics.SuccessfulRuns)
	assert.Equal(t, 34863, insightsSummary.Workflows[0].Metrics.TotalCredits)
}
