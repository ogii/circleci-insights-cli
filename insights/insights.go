package insights

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/jedib0t/go-pretty/list"
	"github.com/jedib0t/go-pretty/table"
	"github.com/ogii/circleci-insights-cli/data"
)

func PrintInsightsSummaryList(insights data.InsightsSummary, dataType string) {
	if !CheckIfWorkflowsDataEmpty(insights) {
		return
	}

	l := list.NewWriter()
	for _, item := range insights.Workflows {
		l.AppendItem("-----------------------------")
		l.AppendItem(fmt.Sprintf("%sName: %s", dataType, item.Name))
		l.AppendItem(fmt.Sprintf("Credits Consumed: %d", item.Metrics.TotalCredits))
		l.AppendItem(fmt.Sprintf("Success Rate: %s", calculateSuccessRate(item.Metrics.SuccessRate)))
		l.AppendItem(fmt.Sprintf("Total Runs: %d", item.Metrics.TotalRuns))
		l.AppendItem(fmt.Sprintf("Failed Runs: %d", item.Metrics.FailedRuns))
		l.AppendItem(fmt.Sprintf("Successful Runs: %d", item.Metrics.SuccessfulRuns))
	}
	fmt.Println(l.Render())
}

func PrintInsightsSummaryTable(insights data.InsightsSummary, dataType string) {
	if !CheckIfWorkflowsDataEmpty(insights) {
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

func OutputInsightsSummaryToCSV(insights data.InsightsSummary, dataType string, path string) error {
	file, err := generateFile(path, 0770)

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

func OutputInsightsSummaryToJSON(insights data.InsightsSummary) error {
	if !CheckIfWorkflowsDataEmpty(insights) {
		return fmt.Errorf("No data available")
	}

	insightsData := []map[string]interface{}{}

	for _, item := range insights.Workflows {
		insightsData = append(insightsData, map[string]interface{}{
			"name":            item.Name,
			"credits_used":    item.Metrics.TotalCredits,
			"successful_runs": item.Metrics.SuccessfulRuns,
			"failed_runs":     item.Metrics.FailedRuns,
			"success_rate":    calculateSuccessRate(item.Metrics.SuccessRate),
		})
	}

	jsonData, err := json.MarshalIndent(insightsData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %v", err)
	}

	fmt.Println(string(jsonData))

	return nil
}

func CheckIfWorkflowsDataEmpty(insights data.InsightsSummary) bool {
	if len(insights.Workflows) == 0 {
		fmt.Println("No data available.")
		return false
	}
	return true
}

func calculateSuccessRate(rate float64) string {
	return fmt.Sprintf("%.3f%%", rate*100)
}

func generateFile(path string, permissions os.FileMode) (*os.File, error) {
	dirPath := filepath.Dir(path)

	err := os.MkdirAll(dirPath, permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	return file, nil
}
