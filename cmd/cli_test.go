package cmd

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitRootCmd(t *testing.T) {
	cmd := InitRootCmd()

	// Verify the root command properties
	assert.Equal(t, "go-gator", cmd.Use, "Command use should be 'go-gator'")
	assert.Equal(t, "Aggregate news from various sources", cmd.Short, "Command short description should match")
	assert.Contains(t, cmd.Long, "Fetch latest and filtered news from different sources", "Command long description should contain 'Fetch latest and filtered news from different sources'")
	assert.Equal(t, "0.0.3", cmd.Version, "Command version should be '0.0.3'")

	// Verify subcommands
	subCmd := cmd.Commands()
	assert.Equal(t, 1, len(subCmd), "There should be one subcommand")

	fetchNewsCmd := subCmd[0]
	assert.Equal(t, "fetch", fetchNewsCmd.Use, "Subcommand use should be 'fetch-news'")
	assert.Equal(t, "Fetching news from downloaded data", fetchNewsCmd.Short, "Subcommand short description should match")
}
