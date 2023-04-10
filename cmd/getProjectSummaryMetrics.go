/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/joho/godotenv"
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

var (
	client   *http.Client
	once     sync.Once
	token    string
	baseURL  string
	loadEnvs = func() {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
		token = os.Getenv("CIRCLECI_TOKEN")
		baseURL = os.Getenv("API_URL")
	}
)

// getProjectSummaryMetricsCmd represents the getProjectSummaryMetrics command
var getProjectSummaryMetricsCmd = &cobra.Command{
	Use:   "getProjectSummaryMetrics",
	Short: "Get summary metrics for a project's workflows",
	Long:  `Get summary metrics for a project's workflows. Workflow runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getProjectSummaryMetrics called")
		client = &http.Client{
			Timeout: 10 * time.Second,
			Transport: &http.Transport{
				MaxIdleConnsPerHost: 5,
			},
		}
	},
}

func init() {
	rootCmd.AddCommand(getProjectSummaryMetricsCmd)
	getProjectSummaryMetricsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectSummaryMetricsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectSummaryMetricsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectSummaryMetricsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
}
