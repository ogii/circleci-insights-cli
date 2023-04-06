/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

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

// getProjectSummaryMetricsCmd represents the getProjectSummaryMetrics command
var getProjectSummaryMetricsCmd = &cobra.Command{
	Use:   "getProjectSummaryMetrics",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getProjectSummaryMetrics called")
	},
}

func init() {
	rootCmd.AddCommand(getProjectSummaryMetricsCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getProjectSummaryMetricsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getProjectSummaryMetricsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
