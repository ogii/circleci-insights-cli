package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/table"
	"github.com/ogii/circleci-insights-cli/data"
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

		insightsSummary, err := fetchInsightsSummary(slug, branch, reportingWindow)
		if err != nil {
			log.Fatal(err)
		}

		switch formatType := format; formatType {
		case "table":
			printInsightsSummaryTable(*insightsSummary)
		default:
			printInsightsSummaryList(*insightsSummary)
		}

	},
}

func fetchInsightsSummary(slug, branch, reportingWindow string) (*data.InsightsSummary, error) {
	url := fmt.Sprintf("%s/insights/%s/workflows/?branch=%s&reporting-window=%s", baseURL, slug, branch, reportingWindow)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %v", err)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 5,
		},
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Circle-Token", token)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error: non-200 response from CircleCI API: %s", resp.Status)
	}

	var insightsSummary data.InsightsSummary
	err = json.NewDecoder(resp.Body).Decode(&insightsSummary)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}
	return &insightsSummary, nil
}

func printInsightsSummaryList(insights data.InsightsSummary) {
	if len(insights.Items) > 0 {
		l := list.NewWriter()
		for _, item := range insights.Items {
			l.AppendItem("-----------------------------")
			l.AppendItem(fmt.Sprintf("Workflow Name: %s", item.Name))
			l.AppendItem(fmt.Sprintf("Credits Consumed: %d", item.Metrics.TotalCredits))
			l.AppendItem(fmt.Sprintf("Success Rate: %.2f%%", item.Metrics.SuccessRate*100))
			l.AppendItem(fmt.Sprintf("Total Runs: %d", item.Metrics.TotalRuns))
			l.AppendItem(fmt.Sprintf("Failed Runs: %d", item.Metrics.FailedRuns))
			l.AppendItem(fmt.Sprintf("Successful Runs: %d", item.Metrics.SuccessfulRuns))
		}
		fmt.Println(l.Render())
	} else {
		fmt.Println("No data available.")
	}
}

func printInsightsSummaryTable(insights data.InsightsSummary) {
	itemsCount := len(insights.Items)
	if itemsCount == 0 {
		fmt.Println("No data available.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Workflow", "Credits Consumed", "Successful Runs", "Failed Runs", "Success Rate"})

	for _, item := range insights.Items {
		t.AppendRow(table.Row{item.Name, item.Metrics.TotalCredits, item.Metrics.SuccessfulRuns, item.Metrics.FailedRuns, item.Metrics.SuccessRate * 100})
	}

	t.Render()
}

func init() {
	rootCmd.AddCommand(getProjectSummaryMetricsCmd)
	getProjectSummaryMetricsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectSummaryMetricsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectSummaryMetricsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectSummaryMetricsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
	getProjectSummaryMetricsCmd.MarkFlagRequired("slug")
}
