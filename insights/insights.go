package insights

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/table"
	"github.com/ogii/circleci-insights-cli/data"
)

func PrintInsightsSummaryList(insights data.InsightsSummary, dataType string) {
	if len(insights.Workflows) > 0 {
		l := list.NewWriter()
		for _, item := range insights.Workflows {
			l.AppendItem("-----------------------------")
			l.AppendItem(fmt.Sprintf("%sName: %s", dataType, item.Name))
			l.AppendItem(fmt.Sprintf("Credits Consumed: %d", item.Metrics.TotalCredits))
			l.AppendItem(fmt.Sprintf("Success Rate: %.3f%%", item.Metrics.SuccessRate*100))
			l.AppendItem(fmt.Sprintf("Total Runs: %d", item.Metrics.TotalRuns))
			l.AppendItem(fmt.Sprintf("Failed Runs: %d", item.Metrics.FailedRuns))
			l.AppendItem(fmt.Sprintf("Successful Runs: %d", item.Metrics.SuccessfulRuns))
		}
		fmt.Println(l.Render())
	} else {
		fmt.Println("No data available.")
	}
}

func PrintInsightsSummaryTable(insights data.InsightsSummary, dataType string) {
	itemsCount := len(insights.Workflows)
	if itemsCount == 0 {
		fmt.Println("No data available.")
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	if dataType == "workflow" {
		t.SetTitle("Workflows")
	} else if dataType == "job" {
		t.SetTitle("Jobs")
	}

	t.AppendHeader(table.Row{"Name", "Credits Consumed", "Successful Runs", "Failed Runs", "Success Rate"})

	for _, item := range insights.Workflows {
		t.AppendRow(table.Row{item.Name, item.Metrics.TotalCredits, item.Metrics.SuccessfulRuns, item.Metrics.FailedRuns, fmt.Sprintf("%.3f%%", item.Metrics.SuccessRate*100)})
	}
	t.SetStyle(table.StyleBold)
	t.Style().Options.SeparateRows = true
	t.Render()
}

func OutputInsightsSummaryToCSV(insights data.InsightsSummary, dataType string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"", "", dataType, "", ""})

	writer.Write([]string{"Name", "Credits Consumed", "Successful Runs", "Failed Runs", "Success Rate"})

	for _, item := range insights.Workflows {
		successRate := fmt.Sprintf("%.3f%%", item.Metrics.SuccessRate*100)
		err = writer.Write([]string{item.Name, strconv.Itoa(item.Metrics.TotalCredits), strconv.Itoa(item.Metrics.SuccessfulRuns), strconv.Itoa(item.Metrics.FailedRuns), successRate})
		if err != nil {
			return fmt.Errorf("failed to write to file: %v", err)
		}
	}

	if err := writer.Error(); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}
	return nil
}
