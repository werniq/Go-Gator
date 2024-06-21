package cmd

import (
	"github.com/spf13/cobra"
	"log"
)

// InitRootCmd initializes root cmd and attaches fetchNews command to our main command
func InitRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "go-gator",
		Short: "Aggregate news from various sources",
		Long: "Fetch latest and filtered news from different sources: NYT, BBC, ABC, etc. \n " +
			"Filter them by topic, key words, country, and timestamp",
		Version: "0.0.3",

		Run: func(cmd *cobra.Command, args []string) {
			log.Println("[Go Gator]")
		},
	}
	rootCmd.AddCommand(AddFetchNewsCmd())

	return rootCmd
}
