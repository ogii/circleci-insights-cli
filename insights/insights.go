package insights

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/table"
	"github.com/ogii/circleci-insights-cli/data"
)

func PrintInsightsSummaryList(insights data.InsightsSummary) {
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

func PrintInsightsSummaryTable(insights data.InsightsSummary) {
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
	t.SetStyle(table.StyleBold)
	t.Style().Options.SeparateRows = true
	t.Render()
}
