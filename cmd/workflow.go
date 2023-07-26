package cmd

import (
	"fmt"
	"log"
	"path"
	"time"

	"github.com/ogii/circleci-insights-cli/client"
	"github.com/ogii/circleci-insights-cli/insights"
	"github.com/spf13/cobra"
)

var getProjectWorkflowJobsCmd = &cobra.Command{
	Use:   "workflow-job-summary",
	Short: "Get summary metrics for a project workflow's jobs.",
	Long:  `Get summary metrics for a project workflow's jobs. Job runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		slug, _ := cmd.Flags().GetString("slug")
		branch, _ := cmd.Flags().GetString("branch")
		format, _ := cmd.Flags().GetString("format")
		reportingWindow, _ := cmd.Flags().GetString("reporting-window")
		workflow, _ := cmd.Flags().GetString("workflow")
		output, _ := cmd.Flags().GetString("output")
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
			var fullpath = path.Join(output, time.Now().UTC().Format("2006-01-02T15_04_05Z")+"_output.csv")
			println(fullpath)
			insights.OutputInsightsSummaryToCSV(*insightsSummary, "Jobs", fullpath)
		case "json":
			insights.OutputInsightsSummaryToJSON(*insightsSummary)
		default:
			insights.PrintInsightsSummaryList(*insightsSummary, "Jobs")
		}
	},
}

func init() {
	getCmd.AddCommand(getProjectWorkflowJobsCmd)
	getProjectWorkflowJobsCmd.Flags().String("slug", "", "The slug for a CircleCI project in the format `gh/orgname/projectname`")
	getProjectWorkflowJobsCmd.Flags().String("branch", "main", "The branch for a CircleCI project")
	getProjectWorkflowJobsCmd.Flags().String("format", "table", "The format of the results to be shown")
	getProjectWorkflowJobsCmd.Flags().String("reporting-window", "last-90-days", "The time window used to calculate summary metrics.")
	getProjectWorkflowJobsCmd.Flags().String("workflow", "", "The name of the workflow")
	getProjectWorkflowJobsCmd.Flags().String("output", "", "The location to save insights api output.")
	getProjectWorkflowJobsCmd.MarkFlagRequired("slug")
	getProjectWorkflowJobsCmd.MarkFlagRequired("workflow")
}
