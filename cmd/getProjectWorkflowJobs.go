package cmd

import (
	"fmt"
	"log"

	"github.com/ogii/circleci-insights-cli/client"
	"github.com/ogii/circleci-insights-cli/insights"
	"github.com/spf13/cobra"
)

var getProjectWorkflowJobsCmd = &cobra.Command{
	Use:   "getProjectWorkflowJobs",
	Short: "Get summary metrics for a project workflow's jobs.",
	Long:  `Get summary metrics for a project workflow's jobs. Job runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, _ := cmd.Flags().GetString("slug")
		branch, _ := cmd.Flags().GetString("branch")
		format, _ := cmd.Flags().GetString("format")
		reportingWindow, _ := cmd.Flags().GetString("reporting-window")
		workflow, _ := cmd.Flags().GetString("workflow")
		url := fmt.Sprintf("%s/jobs", workflow)

		client := client.NewClient(baseURL, token)
		insightsSummary, err := client.FetchInsightsSummary(slug, url, branch, reportingWindow)
		if err != nil {
			log.Fatal(err)
		}

		switch formatType := format; formatType {
		case "table":
			insights.PrintInsightsSummaryTable(*insightsSummary, "Jobs")
		case "csv":
			insights.OutputInsightsSummaryToCSV(*insightsSummary, "Jobs", "output.csv")
		default:
			insights.PrintInsightsSummaryList(*insightsSummary, "Jobs")
		}
	},
}

func init() {
	rootCmd.AddCommand(getProjectWorkflowJobsCmd)
	getProjectWorkflowJobsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectWorkflowJobsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectWorkflowJobsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectWorkflowJobsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
	getProjectWorkflowJobsCmd.Flags().String("workflow", "", "The name of the workflow")
	getProjectWorkflowJobsCmd.MarkFlagRequired("slug")
	getProjectWorkflowJobsCmd.MarkFlagRequired("workflow")
}
