package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ogii/circleci-insights-cli/data"
	"github.com/ogii/circleci-insights-cli/insights"
	"github.com/spf13/cobra"
)

var getProjectSummaryMetricsCmd = &cobra.Command{
	Use:   "getProjectSummaryMetrics",
	Short: "Get summary metrics for a project's workflows",
	Long:  `Get summary metrics for a project's workflows. Workflow runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, _ := cmd.Flags().GetString("slug")
		branch, _ := cmd.Flags().GetString("branch")
		format, _ := cmd.Flags().GetString("format")
		reportingWindow, _ := cmd.Flags().GetString("reporting-window")

		insightsSummary, err := fetchInsightsSummary(client, slug, branch, reportingWindow)
		if err != nil {
			log.Fatal(err)
		}

		switch formatType := format; formatType {
		case "table":
			insights.PrintInsightsSummaryTable(*insightsSummary)
		default:
			insights.PrintInsightsSummaryList(*insightsSummary)
		}

	},
}

var client = &http.Client{
	Timeout: 10 * time.Second,
	Transport: &http.Transport{
		MaxIdleConnsPerHost: 5,
	},
}

func fetchInsightsSummary(client *http.Client, slug, branch, reportingWindow string) (*data.InsightsSummary, error) {
	var insightsSummary data.InsightsSummary
	var nextPageToken string

	for {
		url := fmt.Sprintf("%s/insights/%s/workflows/?branch=%s&reporting-window=%s", baseURL, slug, branch, reportingWindow)
		if nextPageToken != "" {
			url += "&page-token=" + nextPageToken
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating HTTP request for URL %s: %v", url, err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Circle-Token", token)

		resp, err := client.Do(req)
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

func init() {
	rootCmd.AddCommand(getProjectSummaryMetricsCmd)
	getProjectSummaryMetricsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectSummaryMetricsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectSummaryMetricsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectSummaryMetricsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
	getProjectSummaryMetricsCmd.MarkFlagRequired("slug")
}
