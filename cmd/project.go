package cmd

import (
	"log"

	"github.com/ogii/circleci-insights-cli/client"
	"github.com/ogii/circleci-insights-cli/insights"
	"github.com/spf13/cobra"
)

var getProjectSummaryMetricsCmd = &cobra.Command{
	Use:   "project-workflow-summary",
	Short: "Get summary metrics for a project's workflows",
	Long:  `Get summary metrics for a project's workflows. Workflow runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, _ := cmd.Flags().GetString("slug")
		branch, _ := cmd.Flags().GetString("branch")
		format, _ := cmd.Flags().GetString("format")
		reportingWindow, _ := cmd.Flags().GetString("reporting-window")
		url := ""

		client := client.NewClient(baseURL, token)
		insightsSummary, err := client.FetchInsightsSummary(slug, url, branch, reportingWindow)
		if err != nil {
			log.Fatal(err)
		}

		switch formatType := format; formatType {
		case "table":
			insights.PrintInsightsSummaryTable(*insightsSummary, "Workflows")
		case "csv":
			insights.OutputInsightsSummaryToCSV(*insightsSummary, "Workflows", "output.csv")
		case "json":
			insights.OutputInsightsSummaryToJSON(*insightsSummary)
		default:
			insights.PrintInsightsSummaryList(*insightsSummary, "Workflows")
		}

	},
}

func init() {
	getCmd.AddCommand(getProjectSummaryMetricsCmd)
	getProjectSummaryMetricsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectSummaryMetricsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectSummaryMetricsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectSummaryMetricsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
	getProjectSummaryMetricsCmd.MarkFlagRequired("slug")
}
