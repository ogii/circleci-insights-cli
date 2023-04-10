package cmd

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

var (
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
	rootCmd = &cobra.Command{
		Use:   "circleci-insights-cli",
		Short: "A cli tool to get CircleCI insights data",
		Long:  `A cli tool to get CircleCI insights data`,
	}
)

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	loadEnvs()
}
