package cmd

import (
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
	"newsAggr/logger"
	"testing"
)

func TestInitRootCmd(t *testing.T) {
	testRootCmd := &cobra.Command{
		Use:   "go-gator",
		Short: "Aggregate news from various sources",
		Long: "Fetch latest and filtered news from different sources: NYT, BBC, ABC, etc. \n " +
			"Filter them by topic, key words, country, and timestamp",
		Version: "0.0.3",

		Run: func(cmd *cobra.Command, args []string) {
			logger.InfoLogger.Println("[Go Gator]")
		},
	}

	testRootCmd.AddCommand(addFetchNewsCmd())

	assert.Equal(t, testRootCmd.Use, rootCmd.Use)
	assert.Equal(t, testRootCmd.Short, rootCmd.Short)
}
