package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/table"
	"github.com/spf13/cobra"
)

type InsightsSummary struct {
	Items Items  `json:"items"`
	ID    string `json:"id"`
	Name  string `json:"name"`
}

type Items []struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	WindowStart string     `json:"window_start"`
	WindowEnd   string     `json:"window_end"`
	Repository  Repository `json:"repository"`
	Metrics     Metrics    `json:"metrics"`
}

type Repository struct {
	VcsType   string `json:"vcs_type"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type Metrics struct {
	SuccessRate    float64 `json:"success_rate"`
	TotalRuns      int     `json:"total_runs"`
	FailedRuns     int     `json:"failed_runs"`
	SuccessfulRuns int     `json:"successful_runs"`
	TotalCredits   int     `json:"total_credits_used"`
}

var getProjectSummaryMetricsCmd = &cobra.Command{
	Use:   "getProjectSummaryMetrics",
	Short: "Get summary metrics for a project's workflows",
	Long:  `Get summary metrics for a project's workflows. Workflow runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getProjectSummaryMetrics called")

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

func fetchInsightsSummary(slug, branch, reportingWindow string) (*InsightsSummary, error) {
	if slug == "" || branch == "" {
		return nil, errors.New("slug and branch must not be empty")
	}

	fmt.Println(baseURL)
	url := fmt.Sprintf("%s/insights/%s/workflows/?branch=%s&reporting-window=%s", baseURL, slug, branch, reportingWindow)
	fmt.Println(url)
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

	var insightsSummary InsightsSummary
	err = json.NewDecoder(resp.Body).Decode(&insightsSummary)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling JSON response: %v", err)
	}
	return &insightsSummary, nil
}

func printInsightsSummaryList(insights InsightsSummary) {
	if len(insights.Items) > 0 {
		for _, item := range insights.Items {
			fmt.Println("-----------------------------")
			fmt.Printf("Workflow Name: %s\n", item.Name)
			fmt.Printf("Credits Consumed: %d\n", item.Metrics.TotalCredits)
			fmt.Printf("Success Rate: %.2f%%\n", item.Metrics.SuccessRate*100)
			fmt.Printf("Total Runs: %d\n", item.Metrics.TotalRuns)
			fmt.Printf("Failed Runs: %d\n", item.Metrics.FailedRuns)
			fmt.Printf("Successful Runs: %d\n", item.Metrics.SuccessfulRuns)
		}
	} else {
		fmt.Println("No data available.")
	}
}

func printInsightsSummaryTable(insights InsightsSummary) {
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
}
