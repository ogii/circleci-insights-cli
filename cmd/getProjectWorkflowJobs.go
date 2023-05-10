/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// getProjectWorkflowJobsCmd represents the getProjectWorkflowJobs command
var getProjectWorkflowJobsCmd = &cobra.Command{
	Use:   "getProjectWorkflowJobs",
	Short: "Get summary metrics for a project workflow's jobs.",
	Long:  `Get summary metrics for a project workflow's jobs. Job runs going back at most 90 days are included in the aggregation window.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("getProjectWorkflowJobs called")
	},
}

func init() {
	rootCmd.AddCommand(getProjectWorkflowJobsCmd)
}
