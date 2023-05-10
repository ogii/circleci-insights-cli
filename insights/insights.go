package insights

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/table"
	"github.com/ogii/circleci-insights-cli/data"
)

func PrintInsightsSummaryList(insights data.InsightsSummary) {
	if len(insights.Workflows) > 0 {
		l := list.NewWriter()

		for _, workflow := range insights.Workflows {
			l.AppendItem("-----------------------------")
			l.AppendItem(fmt.Sprintf("Workflow Name: %s", workflow.Name))
			l.AppendItem(fmt.Sprintf("Credits Consumed: %d", workflow.Metrics.TotalCredits))
			l.AppendItem(fmt.Sprintf("Success Rate: %.2f%%", workflow.Metrics.SuccessRate*100))
			l.AppendItem(fmt.Sprintf("Total Runs: %d", workflow.Metrics.TotalRuns))
			l.AppendItem(fmt.Sprintf("Failed Runs: %d", workflow.Metrics.FailedRuns))
			l.AppendItem(fmt.Sprintf("Successful Runs: %d", workflow.Metrics.SuccessfulRuns))
		}
		fmt.Println(l.Render())
	} else {
		fmt.Println("No data available.")
	}
}

func PrintInsightsSummaryTable(insights data.InsightsSummary) {
	if len(insights.Workflows) > 0 {
		t := table.NewWriter()
		t.SetOutputMirror(os.Stdout)
		t.AppendHeader(table.Row{"Workflow", "Credits Consumed", "Successful Runs", "Failed Runs", "Success Rate"})

		for _, item := range insights.Workflows {
			t.AppendRow(table.Row{item.Name, item.Metrics.TotalCredits, item.Metrics.SuccessfulRuns, item.Metrics.FailedRuns, item.Metrics.SuccessRate * 100})
		}
		t.Render()
	} else {
		fmt.Println("No data available.")
	}
}
