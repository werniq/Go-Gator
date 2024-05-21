package cmd

import (
	"github.com/spf13/cobra"
	"newsAggr/logger"
)

var rootCmd = &cobra.Command{
	Use:   "go-gator",
	Short: "Aggregate news from various sources",
	Long: "Fetch latest and filtered news from different sources: NYT, BBC, ABC, etc. \n " +
		"Filter them by topic, key words, country, and timestamp",
	Version: "0.0.1",

	Run: func(cmd *cobra.Command, args []string) {
		logger.InfoLogger.Println("Example of first command")
	},
}

func InitRootCmd() *cobra.Command {
	rootCmd.AddCommand(addFetchNewsCmd())

	return rootCmd
}
